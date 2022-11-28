package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func parseGlobalNode(cfg *MachConfig, globalNode *yaml.Node) error {
	if err := globalNode.Decode(&cfg.Global); err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}

	knownKeys := []string{"cloud", "terraform_config", "environment"}
	nodes := mapYamlNodes(globalNode.Content)

	err := iterateYamlNodes(nodes, knownKeys, func(key string, data map[string]any) error {
		return cfg.Plugins.SetGlobalConfig(key, data)
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
				return err
			}
			cfg.Global.TerraformStateProvider = "aws"
		}
		if child, ok := childs["azure_remote_state"]; ok {
			data, err := nodeAsMap(child)
			if err != nil {
				return err
			}
			if err := cfg.Plugins.SetRemoteState("azure", data); err != nil {
				return err
			}
			cfg.Global.TerraformStateProvider = "azure"
		}
	}

	return nil
}
