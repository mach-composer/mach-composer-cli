package config

import (
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"gopkg.in/yaml.v3"
)

func parseGlobalNode(cfg *MachConfig, globalNode *yaml.Node) error {
	if err := globalNode.Decode(&cfg.Global); err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}

	for _, plugin := range cfg.Plugins.All() {
		provider := cfg.Global.TerraformConfig.Providers[plugin.Name]
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
		children := mapYamlNodes(node.Content)

		// Backwards compat
		if child, ok := children["aws_remote_state"]; ok {
			cli.DeprecationWarning(&cli.DeprecationOptions{
				Message: "the usage of `aws_remote_state` is deprecated and will be removed in the next major version",
				Details: `
				Please move the configuration to the remote_state block and add the provider name as plugin.
				
				For example:
				
				    aws_remote_state:
					  key_prefix: mach-composer
					  region: eu-central-1
					  bucket: "mcc-terraform-state"
				
				To:
				
				    remote_state:
					  plugin: aws
					  key_prefix: mach-composer
					  region: eu-central-1
					  bucket: "mcc-terraform-state"
				`,
			})
			data, err := nodeAsMap(child)
			if err != nil {
				return err
			}

			cfg.Global.TerraformConfig.RemoteState = data
			cfg.Global.TerraformStateProvider = "aws"
			return nil
		} else if child, ok := children["azure_remote_state"]; ok {
			cli.DeprecationWarning(&cli.DeprecationOptions{
				Message: "the usage of `azure_remote_state` is deprecated and will be removed in the next major version",
				Details: `
				Please move the configuration to the remote_state block and add the provider name as plugin.
				
				For example:
				
				    azure_remote_state:
						resource_group: some-resource-group
						storage_account: some-account
						container_name: some-container
				
				To:
				
				    remote_state:
						plugin: azure
						resource_group: some-resource-group
						storage_account: some-account
						container_name: some-container
				`,
			})
			data, err := nodeAsMap(child)
			if err != nil {
				return err
			}

			cfg.Global.TerraformConfig.RemoteState = data
			cfg.Global.TerraformStateProvider = "azure"
			return nil
		} else if child, ok := children["remote_state"]; ok {
			data, err := nodeAsMap(child)
			if err != nil {
				return err
			}

			pluginName, ok := data["plugin"].(string)
			if !ok {
				return fmt.Errorf("plugin needs to be defined for remote_state")
			}

			cfg.Global.TerraformStateProvider = pluginName
			return nil
		}
	}

	return nil
}
