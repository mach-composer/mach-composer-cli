package updater

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/labd/mach-composer/internal/config"
)

type WorkerJob struct {
	component *config.Component
	cfg       *config.MachConfig
}

type Updater struct {
	filename string
	config   *config.MachConfig
	updates  []ChangeSet
}

func NewUpdater(ctx context.Context, filename string) (*Updater, error) {
	cfg, err := config.Load(ctx, filename, "")
	if err != nil {
		return nil, err
	}

	result := &Updater{
		filename: filename,
		config:   cfg,
	}

	return result, nil
}

// UpdateAllComponents updates all the components in the config file.
func (u *Updater) UpdateAllComponents(ctx context.Context) {
	updateSet := FindUpdates(ctx, u.config, u.filename)
	u.updates = updateSet.updates
	fmt.Printf("%d components have updates available\n", len(u.updates))
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
	if u.config.IsEncrypted {
		SopsFileWriter(ctx, u.config, updateSet)
	} else {
		MachFileWriter(updateSet)
	}

	return true
}

func FindUpdates(ctx context.Context, cfg *config.MachConfig, filename string) *UpdateSet {
	numUpdates := len(cfg.Components)
	jobs := make(chan WorkerJob, numUpdates)
	results := make(chan *ChangeSet, numUpdates)

	fmt.Printf("Checking if there are updates for %d components\n", numUpdates)

	const numWorkers = 4

	// Start 4 workers
	for i := 0; i < numWorkers; i++ {
		go func() {
			for j := range jobs {
				cs, err := GetLastVersion(ctx, j.component, j.cfg.Filename)
				if err != nil {
					panic(err)
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

	// Process results as we receive them from the channel
	updates := UpdateSet{
		filename: filename,
	}
	for i := 0; i < numUpdates; i++ {
		changeSet := <-results

		output := OutputChanges(changeSet)
		fmt.Print(output)

		if changeSet.HasChanges() {
			updates.updates = append(updates.updates, *changeSet)
		}
	}

	return &updates
}

func FindSpecificUpdate(ctx context.Context, cfg *config.MachConfig, filename string, component *config.Component) (*UpdateSet, error) {
	changeSet, err := GetLastVersion(ctx, component, cfg.Filename)
	if err != nil {
		return nil, err
	}

	output := OutputChanges(changeSet)
	fmt.Print(output)

	updates := UpdateSet{
		filename: cfg.Filename,
		updates:  []ChangeSet{*changeSet},
	}
	return &updates, nil
}

func GetLastVersion(ctx context.Context, c *config.Component, origin string) (*ChangeSet, error) {
	if strings.HasPrefix(c.Source, "git:") {
		return GetLastVersionGit(ctx, c, origin)
	}
	return nil, errors.New("unrecognized component source")
}
