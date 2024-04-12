package updater

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"gopkg.in/yaml.v3"

	"github.com/mach-composer/mach-composer-cli/internal/cloud"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type PartialConfig struct {
	MachComposer   config.MachComposer      `yaml:"mach_composer"`
	Components     []config.ComponentConfig `yaml:"components"`
	ComponentsNode *yaml.Node
	Sops           yaml.Node `yaml:"sops"`

	isEncrypted bool
	filename    string `yaml:"-"`
	client      *mccsdk.APIClient
}

type PartialRawConfig struct {
	MachComposer config.MachComposer `yaml:"mach_composer"`
	Components   yaml.Node           `yaml:"components"`
	Sops         yaml.Node           `yaml:"sops"`
}

func (c *PartialConfig) GetComponent(name string) *config.ComponentConfig {
	for i := range c.Components {
		if strings.EqualFold(c.Components[i].Name, name) {
			return &c.Components[i]
		}
	}
	return nil
}

type WorkerJob struct {
	component *config.ComponentConfig
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
	filename string // Should always point to the filename containg the components
	config   *PartialConfig
	updates  []ChangeSet
}

// NewUpdater creates an update to update the component versions in a config
// file.
func NewUpdater(ctx context.Context, filename string, useCloud bool) (*Updater, error) {
	body, err := utils.AFS.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	raw := &PartialRawConfig{}
	err = yaml.Unmarshal(body, raw)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall yaml: %w", err)
	}

	cwd := path.Dir(filename)

	// Resolve $ref and include() references for the components.
	componentNode, componentsFilename, err := config.LoadRefData(ctx, &raw.Components, cwd)
	if err != nil {
		return nil, err
	}
	raw.Components = *componentNode

	var components []config.ComponentConfig
	if err := raw.Components.Decode(&components); err != nil {
		return nil, fmt.Errorf("decoding error: %w", err)
	}

	cfg := &PartialConfig{
		MachComposer:   raw.MachComposer,
		Components:     components,
		ComponentsNode: &raw.Components,
		Sops:           raw.Sops,
		filename:       filepath.Base(filename),
	}

	// If we have a Sops node which is a mapping then we can assume that this
	// file is encrypted. This is only relevant if the components are stored in
	// main config file, and not referenced
	if componentsFilename == "" {
		cfg.isEncrypted = cfg.Sops.Kind == yaml.MappingNode
	}

	result := &Updater{
		filename: filename,
		config:   cfg,
	}

	if componentsFilename != "" {
		result.filename = path.Join(path.Dir(filename), componentsFilename)
	}

	if useCloud {
		if cfg.MachComposer.Cloud.Empty() {
			return nil, fmt.Errorf("please defined cloud details")
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
		sopsFileWriter(ctx, u.config, updateSet)
	} else {
		machFileWriter(ctx, u.config, updateSet)
	}

	return true
}
