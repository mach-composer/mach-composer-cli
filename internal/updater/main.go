package updater

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/labd/mach-composer/internal/config"
	"github.com/labd/mach-composer/internal/utils"
)

type PartialConfig struct {
	Components []config.Component `yaml:"components"`
	Sops       yaml.Node          `yaml:"sops"`

	isEncrypted bool
	filename    string `yaml:"-"`
}

func (c *PartialConfig) GetComponent(name string) *config.Component {
	for i := range c.Components {
		if strings.EqualFold(c.Components[i].Name, name) {
			return &c.Components[i]
		}
	}
	return nil
}

type PartialRawConfig struct {
	Components yaml.Node
}

type WorkerJob struct {
	component *config.Component
	cfg       *PartialConfig
}

type UpdateError struct {
	msg       string
	component string
	source    string
}

func (e *UpdateError) Error() string {
	return e.msg
}

type Updater struct {
	filename string
	config   *PartialConfig
	updates  []ChangeSet
}

func NewUpdater(ctx context.Context, filename string) (*Updater, error) {
	body, err := utils.AFS.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg := &PartialConfig{}
	err = yaml.Unmarshal(body, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall yaml: %w", err)
	}
	cfg.filename = filepath.Base(filename)

	// If we have a Sops node which is a mapping then we can assume that this
	// file is encrypted.
	cfg.isEncrypted = cfg.Sops.Kind == yaml.MappingNode

	result := &Updater{
		filename: filename,
		config:   cfg,
	}

	return result, nil
}

// UpdateAllComponents updates all the components in the config file.
func (u *Updater) UpdateAllComponents(ctx context.Context) error {
	updateSet, err := FindUpdates(ctx, u.config, u.filename)
	if err != nil {
		return err
	}
	u.updates = updateSet.updates
	fmt.Printf("%d components have updates available\n", len(u.updates))
	return nil
}

// UpdateComponent updates a specific component. When the version is empty it
// will retrieve the last version for the given component.
func (u *Updater) UpdateComponent(ctx context.Context, name, version string) error {
	component := u.config.GetComponent(name)
	if component == nil {
		return fmt.Errorf("No component found with name %s", name)
	}

	// If no specific version is defined we auto-detect the last version
	if version != "" {
		u.updates = []ChangeSet{
			{
				LastVersion: version,
				Component:   component,
				Forced:      true,
			},
		}
		fmt.Printf("Setting component %s to version %s\n", component.Name, version)
		return nil
	}

	updateSet, err := FindSpecificUpdate(ctx, u.config, u.filename, component)
	if err != nil {
		return err
	}
	if updateSet.HasChanges() {
		fmt.Printf("Updating component %s to version %s\n", component.Name, updateSet.updates[0].LastVersion)
		u.updates = updateSet.updates
	} else {
		fmt.Printf("No updates for component %s\n", component.Name)
	}
	return nil
}

func (u *Updater) GetUpdateSet() *UpdateSet {
	return &UpdateSet{
		filename: u.filename,
		updates:  u.updates,
	}
}

func (u *Updater) Write(ctx context.Context) bool {
	if u.updates == nil || len(u.updates) < 1 {
		return false
	}

	updateSet := u.GetUpdateSet()
	if u.config.isEncrypted {
		SopsFileWriter(ctx, u.config, updateSet)
	} else {
		MachFileWriter(updateSet)
	}

	return true
}

func FindUpdates(ctx context.Context, cfg *PartialConfig, filename string) (*UpdateSet, error) {
	numUpdates := len(cfg.Components)
	jobs := make(chan WorkerJob, numUpdates)
	results := make(chan *ChangeSet, numUpdates)
	errors := make(chan error, numUpdates)

	fmt.Printf("Checking if there are updates for %d components\n", numUpdates)

	var wg sync.WaitGroup
	const numWorkers = 4

	// Start 4 workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := range jobs {
				cs, err := GetLastVersion(ctx, j.component, j.cfg.filename)
				if err != nil {
					errors <- err
					return
				}

				results <- cs
			}
		}()
	}

	// Send work
	for i := range cfg.Components {
		jobs <- WorkerJob{
			component: &cfg.Components[i],
			cfg:       cfg,
		}
	}
	close(jobs)

	wg.Wait()
	close(errors)

	for err := range errors {
		return nil, err
	}

	// Process results as we receive them from the channel
	updates := UpdateSet{
		filename: filename,
	}
	for i := 0; i < numUpdates; i++ {
		changeSet := <-results

		if changeSet == nil {
			continue
		}

		output := OutputChanges(changeSet)
		fmt.Print(output)

		if changeSet.HasChanges() {
			updates.updates = append(updates.updates, *changeSet)
		}
	}

	return &updates, nil
}

func FindSpecificUpdate(ctx context.Context, cfg *PartialConfig, filename string, component *config.Component) (*UpdateSet, error) {
	changeSet, err := GetLastVersion(ctx, component, cfg.filename)
	if err != nil {
		return nil, err
	}

	output := OutputChanges(changeSet)
	fmt.Print(output)

	updates := UpdateSet{
		filename: cfg.filename,
		updates:  []ChangeSet{*changeSet},
	}
	return &updates, nil
}

func GetLastVersion(ctx context.Context, c *config.Component, origin string) (*ChangeSet, error) {
	if strings.HasPrefix(c.Source, "git:") {
		return GetLastVersionGit(ctx, c, origin)
	}

	err := &UpdateError{
		msg:       fmt.Sprintf("unrecognized component source for %s: %s", c.Name, c.Source),
		component: c.Name,
		source:    c.Source,
	}
	return nil, err
}
