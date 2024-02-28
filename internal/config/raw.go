package config

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/mach-composer/mach-composer-cli/internal/plugins"
	"gopkg.in/yaml.v3"
)

type rawConfig struct {
	MachComposer MachComposer `yaml:"mach_composer"`
	Global       yaml.Node    `yaml:"global"`
	Sites        yaml.Node    `yaml:"sites"`
	Components   yaml.Node    `yaml:"components"`

	document  *yaml.Node                `yaml:"-"`
	filename  string                    `yaml:"-"`
	plugins   *plugins.PluginRepository `yaml:"-"`
	variables *Variables                `yaml:"-"`
}

func (r *rawConfig) validate() error {
	if r.MachComposer.Version == "" {
		return fmt.Errorf("no version")
	}

	if r.filename == "" {
		return fmt.Errorf("filename must be set")
	}
	if r.variables == nil {
		return fmt.Errorf("variables cannot be nil")
	}
	if r.plugins == nil {
		return fmt.Errorf("plugins cannot be nil")
	}

	return nil
}

func newRawConfig(filename string, document *yaml.Node) (*rawConfig, error) {
	r := &rawConfig{
		filename:  filename,
		variables: NewVariables(),
		document:  document,
	}

	if err := defaults.Set(r); err != nil {
		return nil, err
	}

	if err := document.Decode(r); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}
	return r, nil
}
