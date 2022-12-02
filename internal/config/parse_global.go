package config

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/labd/mach-composer/internal/plugins"
)

func parseGlobalNode(cfg *MachConfig, globalNode *yaml.Node) error {
	if err := globalNode.Decode(&cfg.Global); err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}

	for identifier, plugin := range cfg.Plugins.All() {
		provider := cfg.Global.TerraformConfig.Providers[identifier]
		if err := plugin.Configure(cfg.Global.Environment, provider); err != nil {
			return err
		}
	}

	knownKeys := []string{"cloud", "terraform_config", "environment"}
	nodes := mapYamlNodes(globalNode.Content)

	err := iterateYamlNodes(nodes, knownKeys, func(pluginName string, data map[string]any) error {
		err := cfg.Plugins.SetGlobalConfig(pluginName, data)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	if node, ok := nodes["terraform_config"]; ok {
		childs := mapYamlNodes(node.Content)

		// Backwards compat
		if child, ok := childs["aws_remote_state"]; ok {
			data, err := nodeAsMap(child)
			if err != nil {
				return err
			}
			if err := cfg.Plugins.SetRemoteState("aws", data); err != nil {
				if _, ok := err.(*plugins.PluginNotFoundError); ok {
					return errors.New("the aws plugin is required when setting aws_remote_state")
				}
				return err
			}
			cfg.Global.TerraformStateProvider = "aws"
			return nil
		} else if child, ok := childs["azure_remote_state"]; ok {
			data, err := nodeAsMap(child)
			if err != nil {
				return err
			}
			if err := cfg.Plugins.SetRemoteState("azure", data); err != nil {
				return err
			}
			cfg.Global.TerraformStateProvider = "azure"
			return nil
		} else if child, ok := childs["remote_state"]; ok {
			data, err := nodeAsMap(child)
			if err != nil {
				return err
			}

			pluginName, ok := data["plugin"].(string)
			if !ok {
				return fmt.Errorf("plugin needs to be defined for remote_state")
			}

			if err := cfg.Plugins.SetRemoteState(pluginName, data); err != nil {
				return err
			}
			cfg.Global.TerraformStateProvider = pluginName
			return nil

		}
	}

	return nil
}
