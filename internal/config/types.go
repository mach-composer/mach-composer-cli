package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/labd/mach-composer/internal/plugins"
	"github.com/labd/mach-composer/internal/variables"
)

type MachConfig struct {
	Filename     string       `yaml:"-"`
	MachComposer MachComposer `yaml:"mach_composer"`
	Global       GlobalConfig `yaml:"global"`
	Sites        []SiteConfig `yaml:"sites"`
	Components   []Component  `yaml:"components"`

	extraFiles  map[string][]byte         `yaml:"-"`
	Plugins     *plugins.PluginRepository `yaml:"-"`
	Variables   *variables.Variables      `yaml:"-"`
	IsEncrypted bool                      `yaml:"-"`
}

func NewMachConfig() *MachConfig {
	cfg := &MachConfig{}
	cfg.extraFiles = make(map[string][]byte, 0)
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

func (c *MachConfig) addFileToConfig(filename string) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading variables file %s: %w", filename, err)
	}
	filename = filepath.Base(filename)
	c.extraFiles[filename] = b
	return nil
}

func (c *MachConfig) GetFiles() map[string][]byte {
	return c.extraFiles
}

type _RawMachConfig struct {
	Filename     string
	MachComposer MachComposer `yaml:"mach_composer"`
	Global       yaml.Node    `yaml:"global"`
	Sites        yaml.Node    `yaml:"sites"`
	Components   yaml.Node    `yaml:"components"`
}

type MachComposer struct {
	Version       string            `yaml:"version"`
	VariablesFile string            `yaml:"variables_file"`
	Plugins       map[string]string `yaml:"plugins"`
}

type GlobalConfig struct {
	Environment            string           `yaml:"environment"`
	Cloud                  string           `yaml:"cloud"`
	TerraformStateProvider string           `yaml:"-"`
	TerraformConfig        *TerraformConfig `yaml:"terraform_config"`
}

// Site contains all configuration needed for a site.
type SiteConfig struct {
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

type TerraformConfig struct {
	Providers map[string]string `yaml:"providers"`
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
