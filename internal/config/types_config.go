package config

import (
	"strings"

	"gopkg.in/yaml.v3"
)

type MachConfig struct {
	Filename     string
	MachComposer MachComposer `yaml:"mach_composer"`
	Global       Global       `yaml:"global"`
	Sites        []Site       `yaml:"sites"`
	Components   []Component  `yaml:"components"`

	ExtraFiles map[string][]byte

	Variables   *Variables
	IsEncrypted bool
}

func NewMachConfig() *MachConfig {
	cfg := &MachConfig{}
	cfg.ExtraFiles = make(map[string][]byte, 0)
	return cfg
}

func (c *MachConfig) HasSite(ident string) bool {
	for i := range c.Sites {
		if c.Sites[i].Identifier == ident {
			return true
		}
	}
	return false
}

func (c *MachConfig) GetComponent(name string) *Component {
	for i := range c.Components {
		if strings.EqualFold(c.Components[i].Name, name) {
			return &c.Components[i]
		}
	}
	return nil
}

type _RawMachConfig struct {
	Filename     string
	MachComposer MachComposer `yaml:"mach_composer"`
	Global       Global       `yaml:"global"`
	Sites        yaml.Node    `yaml:"sites"`
	Components   yaml.Node    `yaml:"components"`
	Sops         yaml.Node    `yaml:"sops"`
}

type MachComposer struct {
	Version       string `yaml:"version"`
	VariablesFile string `yaml:"variables_file"`
}

type Component struct {
	Name         string
	Source       string
	Version      string `yaml:"version"`
	Branch       string
	Integrations []string
	Endpoints    map[string]string `yaml:"endpoints"`

	Azure *ComponentAzureConfig `yaml:"azure"`
}

// UseVersionReference indicates if the module should be referenced with the
// version.
// This will be mainly used for development purposes when referring to a local
// directory; versioning is not possible, but we should still be able to define
// a version in our component for the actual function deployment itself.
func (c *Component) UseVersionReference() bool {
	return strings.HasPrefix(c.Source, "git")
}
