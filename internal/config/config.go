package config

import (
	"github.com/mach-composer/mach-composer-cli/internal/plugins"
	"github.com/mach-composer/mach-composer-cli/internal/state"
	"github.com/mach-composer/mcc-sdk-go/mccsdk"
)

type MachConfig struct {
	Filename     string            `yaml:"-"`
	MachComposer MachComposer      `yaml:"mach_composer"`
	Global       GlobalConfig      `yaml:"global"`
	Sites        SiteConfigs       `yaml:"sites"`
	Components   []ComponentConfig `yaml:"components"`

	StateRepository *state.Repository

	extraFiles  map[string][]byte         `yaml:"-"`
	Plugins     *plugins.PluginRepository `yaml:"-"`
	Variables   *Variables                `yaml:"-"`
	IsEncrypted bool                      `yaml:"-"`
}

func (c *MachConfig) Close() {
	if c.Plugins != nil {
		c.Plugins.Close()
	}
}

func (c *MachConfig) HasSite(ident string) bool {
	for i := range c.Sites {
		if c.Sites[i].Identifier == ident {
			return true
		}
	}
	return false
}

type MachComposer struct {
	Version       any                         `yaml:"version"`
	VariablesFile string                      `yaml:"variables_file"`
	Plugins       map[string]MachPluginConfig `yaml:"plugins"`
	Cloud         MachComposerCloud           `yaml:"cloud"`
	Deployment    Deployment                  `yaml:"deployment"`
}

func (mc *MachComposer) CloudEnabled() bool {
	return !mc.Cloud.Empty()
}

type MachComposerCloud struct {
	Organization string `yaml:"organization"`
	Project      string `yaml:"project"`

	Client *mccsdk.APIClient
}

func (mcc *MachComposerCloud) Empty() bool {
	if mcc.Organization == "" {
		return true
	}
	if mcc.Project == "" {
		return true
	}
	return false
}

type MachPluginConfig struct {
	Source  string `yaml:"source"`
	Version string `yaml:"version"`
	Replace string `yaml:"replace"`
}
