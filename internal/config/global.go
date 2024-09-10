package config

import (
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"gopkg.in/yaml.v3"
)

type GlobalConfig struct {
	Environment            string           `yaml:"environment"`
	Cloud                  string           `yaml:"cloud"`
	TerraformStateProvider string           `yaml:"-"`
	TerraformConfig        *TerraformConfig `yaml:"terraform_config"`
}

type TerraformConfig struct {
	Providers   map[string]string `yaml:"providers"`
	RemoteState map[string]any    `yaml:"remote_state"`
}

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

	nodes := MapYamlNodes(globalNode.Content)

	for _, plugin := range cfg.Plugins.All() {
		data := map[string]any{}

		pluginNode, ok := nodes[plugin.Name]
		if ok {
			var err error
			data, err = nodeAsMap(pluginNode)
			if err != nil {
				return err
			}
		}

		data = utils.FilterMap(data, []string{"cloud", "terraform_config", "environment"})

		if err := plugin.SetGlobalConfig(data); err != nil {
			return fmt.Errorf("%s.SetGlobalConfig failed: %w", plugin.Name, err)
		}
	}

	if node, ok := nodes["terraform_config"]; ok {
		children := MapYamlNodes(node.Content)

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
