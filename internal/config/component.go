package config

import (
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/mach-composer/mach-composer-plugin-sdk/v2/schema"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type ComponentConfig struct {
	Name         string            `yaml:"name"`
	Source       Source            `yaml:"source"`
	Paths        []string          `yaml:"paths"`
	Version      string            `yaml:"version"`
	Branch       string            `yaml:"branch"`
	Integrations []string          `yaml:"integrations"`
	Endpoints    map[string]string `yaml:"endpoints"`
}

func parseComponentsNode(cfg *MachConfig, node *yaml.Node) error {
	if err := node.Decode(&cfg.Components); err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}

	if err := verifyComponents(cfg); err != nil {
		return fmt.Errorf("verify of components failed: %w", err)
	}

	for _, component := range node.Content {
		for _, plugin := range cfg.Plugins.All() {
			data := map[string]any{}
			nodes := mapYamlNodes(component.Content)
			componentName := nodes["name"].Value
			version := nodes["version"].Value

			pluginNode, ok := nodes[plugin.Name]
			if ok {
				var err error
				data, err = nodeAsMap(pluginNode)
				if err != nil {
					return err
				}
			}

			data = utils.FilterMap(data, []string{
				"name", "source", "version", "branch", "integrations", "endpoints", "paths",
			})

			if err := plugin.SetComponentConfig(componentName, version, data); err != nil {
				return fmt.Errorf("%s.SetComponentConfig failed: %w", plugin.Name, err)
			}
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

		// Make sure the component names are unique. Otherwise, raise an error
		if pie.Contains(seen, c.Name) {
			return fmt.Errorf("component %s is duplicate", c.Name)
		}

		// If the component has no integrations (or now called plugins)
		// specified then set it to the cloud integration
		if cfg.Global.Cloud != "" && len(c.Integrations) < 1 {
			c.Integrations = append(c.Integrations, cfg.Global.Cloud)
		}

		seen = append(seen, c.Name)
	}

	return nil
}
