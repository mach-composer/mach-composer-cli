package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/mach-composer/mach-composer-cli/internal/state"

	"github.com/elliotchance/pie/v2"
	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"gopkg.in/yaml.v3"

	"github.com/mach-composer/mach-composer-cli/internal/plugins"
	"github.com/mach-composer/mach-composer-cli/internal/variables"
)

type rawConfig struct {
	MachComposer MachComposer `yaml:"mach_composer"`
	Global       yaml.Node    `yaml:"global"`
	Sites        yaml.Node    `yaml:"sites"`
	Components   yaml.Node    `yaml:"components"`

	document  *yaml.Node                `yaml:"-"`
	filename  string                    `yaml:"-"`
	plugins   *plugins.PluginRepository `yaml:"-"`
	variables *variables.Variables      `yaml:"-"`
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
		variables: variables.NewVariables(),
		document:  document,
	}
	if err := document.Decode(r); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}
	return r, nil
}

type MachConfig struct {
	Filename     string            `yaml:"-"`
	MachComposer MachComposer      `yaml:"mach_composer"`
	Global       GlobalConfig      `yaml:"global"`
	Sites        []SiteConfig      `yaml:"sites"`
	Components   []ComponentConfig `yaml:"components"`

	StateRepository *state.Repository

	extraFiles  map[string][]byte         `yaml:"-"`
	ConfigHash  string                    `yaml:"-"`
	Plugins     *plugins.PluginRepository `yaml:"-"`
	Variables   *variables.Variables      `yaml:"-"`
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

	Definition *ComponentConfig `yaml:"-"`
}

type ComponentConfig struct {
	Name         string
	Source       string
	Paths        []string `yaml:"paths"`
	Version      string   `yaml:"version"`
	Branch       string
	Integrations []string
	Endpoints    map[string]string `yaml:"endpoints"`
}

type TerraformConfig struct {
	Providers   map[string]string `yaml:"providers"`
	RemoteState map[string]any    `yaml:"remote_state"`
}

func (sc SiteComponent) HasCloudIntegration(g *GlobalConfig) bool {
	if sc.Definition == nil {
		log.Fatalf("ComponentConfig %s was not resolved properly (missing definition)", sc.Name)
	}
	return pie.Contains(sc.Definition.Integrations, g.Cloud)
}

// IsGitSource indicates if the source definition refers to Git.
func (c *ComponentConfig) IsGitSource() bool {
	return strings.HasPrefix(c.Source, "git")
}
