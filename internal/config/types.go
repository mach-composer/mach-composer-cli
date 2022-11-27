package config

import (
	"log"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/labd/mach-composer/internal/plugins"
)

type MachConfig struct {
	Filename     string       `yaml:"-"`
	MachComposer MachComposer `yaml:"mach_composer"`
	Global       Global       `yaml:"global"`
	Sites        []Site       `yaml:"sites"`
	Components   []Component  `yaml:"components"`

	ExtraFiles map[string][]byte
	Plugins    *plugins.PluginRepository `yaml:"-"`

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
	Global       yaml.Node    `yaml:"global"`
	Sites        yaml.Node    `yaml:"sites"`
	Components   yaml.Node    `yaml:"components"`
	Sops         yaml.Node    `yaml:"sops"`
}

type MachComposer struct {
	Version       string            `yaml:"version"`
	VariablesFile string            `yaml:"variables_file"`
	Plugins       map[string]string `yaml:"plugins"`
}

type Global struct {
	Environment            string `yaml:"environment"`
	Cloud                  string `yaml:"cloud"`
	TerraformStateProvider string `yaml:"-"`
}

// Site contains all configuration needed for a site.
type Site struct {
	Name         string
	Identifier   string
	RawEndpoints map[string]any `yaml:"endpoints"`

	Components []SiteComponent `yaml:"components"`
}

type SiteComponent struct {
	Name      string
	Variables map[string]any
	Secrets   map[string]any

	Definition *Component `yaml:"-"`
}

type Component struct {
	Name         string
	Source       string
	Version      string `yaml:"version"`
	Branch       string
	Integrations []string
	Endpoints    map[string]string `yaml:"endpoints"`
}

// TODO: should check if the integration is the cloud provider
func (sc SiteComponent) HasCloudIntegration() bool {
	if sc.Definition == nil {
		log.Fatalf("Component %s was not resolved properly (missing definition)", sc.Name)
	}
	for _, i := range sc.Definition.Integrations {
		if i == "aws" || i == "azure" {
			return true
		}
	}
	return false
}

// UseVersionReference indicates if the module should be referenced with the
// version.
// This will be mainly used for development purposes when referring to a local
// directory; versioning is not possible, but we should still be able to define
// a version in our component for the actual function deployment itself.
func (c *Component) UseVersionReference() bool {
	return strings.HasPrefix(c.Source, "git")
}
