package updater

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"gopkg.in/yaml.v3"

	"github.com/labd/mach-composer/internal/cloud"
	"github.com/labd/mach-composer/internal/config"
	"github.com/labd/mach-composer/internal/utils"
)

type PartialConfig struct {
	MachComposer config.MachComposer `yaml:"mach_composer"`
	Components   []config.Component  `yaml:"components"`
	Sops         yaml.Node           `yaml:"sops"`

	isEncrypted bool
	filename    string `yaml:"-"`
	client      *mccsdk.APIClient
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

func NewUpdater(ctx context.Context, filename string, useCloud bool) (*Updater, error) {
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

	if useCloud {
		if cfg.MachComposer.Cloud.Empty() {
			return nil, fmt.Errorf("Please defined cloud details")
		}

		client, err := cloud.NewClient(ctx)
		if err != nil {
			return nil, err
		}
		cfg.client = client
	}

	return result, nil
}

// UpdateAllComponents updates all the components in the config file.
func (u *Updater) UpdateAllComponents(ctx context.Context) error {
	updateSet, err := findUpdates(ctx, u.config, u.filename)
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

	updateSet, err := findSpecificUpdate(ctx, u.config, u.filename, component)
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
