package config

import (
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/mach-composer/mach-composer-cli/internal/config/variable"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"strings"
)

type SiteComponentConfigs []SiteComponentConfig

func (s *SiteComponentConfigs) Get(name string) (*SiteComponentConfig, error) {
	for _, site := range *s {
		if site.Name == name {
			return &site, nil
		}
	}
	return nil, fmt.Errorf("site component %s not found", name)
}

type SiteComponentConfig struct {
	Name       string                `yaml:"name"`
	Definition *ComponentConfig      `yaml:"-"`
	Variables  variable.VariablesMap `yaml:"variables"`
	Secrets    variable.VariablesMap `yaml:"secrets"`
	Deployment *Deployment           `yaml:"deployment"`

	DependsOn []string `yaml:"depends_on"`
}

func (sc *SiteComponentConfig) Hash() (string, error) {
	return utils.ComputeHash(struct {
		Name      string                `json:"name"`
		Version   string                `json:"version"`
		Variables variable.VariablesMap `json:"variables"`
		Secrets   variable.VariablesMap `json:"secrets"`
	}{
		Name:      sc.Name,
		Version:   sc.Definition.Version,
		Variables: sc.Variables,
		Secrets:   sc.Secrets,
	})
}

type ComponentConfig struct {
	Name         string            `yaml:"name"`
	Source       string            `yaml:"source"`
	Paths        []string          `yaml:"paths"`
	Version      string            `yaml:"version"`
	Branch       string            `yaml:"branch"`
	Integrations []string          `yaml:"integrations"`
	Endpoints    map[string]string `yaml:"endpoints"`
}

type TerraformConfig struct {
	Providers   map[string]string `yaml:"providers"`
	RemoteState map[string]any    `yaml:"remote_state"`
}

func (sc *SiteComponentConfig) HasCloudIntegration(g *GlobalConfig) bool {
	if sc.Definition == nil {
		log.Fatal().Msgf("ComponentConfig %s was not resolved properly (missing definition)", sc.Name)
	}
	return pie.Contains(sc.Definition.Integrations, g.Cloud)
}

// IsGitSource indicates if the source definition refers to Git.
func (c *ComponentConfig) IsGitSource() bool {
	return strings.HasPrefix(c.Source, "git")
}

func parseComponentsNode(cfg *MachConfig, node *yaml.Node) error {
	if err := node.Decode(&cfg.Components); err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}

	if err := verifyComponents(cfg); err != nil {
		return fmt.Errorf("verify of components failed: %w", err)
	}
	knownKeys := []string{
		"name", "source", "version", "branch", "integrations", "endpoints", "paths",
	}
	for _, component := range node.Content {
		nodes := mapYamlNodes(component.Content)
		componentName := nodes["name"].Value
		err := iterateYamlNodes(nodes, knownKeys, func(key string, data map[string]any) error {
			return cfg.Plugins.SetComponentConfig(key, componentName, data)
		})
		if err != nil {
			return err
		}
	}

	if err := registerComponentEndpoints(cfg); err != nil {
		return fmt.Errorf("register of components failed: %w", err)
	}

	return nil
}

func registerComponentEndpoints(cfg *MachConfig) error {
	var cloudPlugin schema.MachComposerPlugin
	if cfg.Global.Cloud != "" {
		var err error
		cloudPlugin, err = cfg.Plugins.Get(cfg.Global.Cloud)
		if err != nil {
			return err
		}
	}

	for i := range cfg.Components {
		c := &cfg.Components[i]
		if cloudPlugin == nil {
			if len(c.Endpoints) > 0 {
				log.Error().Msg("Unable to register component endpoints when no cloud provider is configured")
			}
			continue
		}
		if len(c.Endpoints) > 0 {
			err := cloudPlugin.SetComponentEndpointsConfig(c.Name, c.Endpoints)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Verify the components config and set default values where needed.
func verifyComponents(cfg *MachConfig) error {
	var seen []string
	for i := range cfg.Components {
		c := &cfg.Components[i]

		if c.Integrations == nil {
			c.Integrations = make([]string, 0)
		}

		// Make sure the component names are unique. Otherwise raise an error
		if pie.Contains(seen, c.Name) {
			return fmt.Errorf("component %s is duplicate", c.Name)
		}

		// If the component has no integrations (or now called plugins)
		// specified then set it to the cloud integration
		if cfg.Global.Cloud != "" && len(c.Integrations) < 1 {
			c.Integrations = append(c.Integrations, cfg.Global.Cloud)
		}

		// If the source is a relative locale path then transform it to an
		// absolute path (required for Terraform)
		if strings.HasPrefix(c.Source, ".") {
			if val, err := filepath.Abs(c.Source); err == nil {
				c.Source = val
			} else {
				return err
			}
		}

		seen = append(seen, c.Name)
	}

	return nil
}
