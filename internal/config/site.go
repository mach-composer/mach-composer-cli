package config

import (
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/config/variable"
	"github.com/rs/zerolog/log"

	"github.com/elliotchance/pie/v2"
	"gopkg.in/yaml.v3"

	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type SiteConfigs []SiteConfig

func (s *SiteConfigs) Get(identifier string) (*SiteConfig, error) {
	for _, site := range *s {
		if site.Identifier == identifier {
			return &site, nil
		}
	}
	return nil, fmt.Errorf("site %s not found", identifier)
}

// SiteConfig contains all configuration needed for a site.
type SiteConfig struct {
	Name         string         `yaml:"name"`
	Identifier   string         `yaml:"identifier"`
	Deployment   *Deployment    `yaml:"deployment"`
	RawEndpoints map[string]any `yaml:"endpoints"`

	Variables variable.VariablesMap `yaml:"variables"`
	Secrets   variable.VariablesMap `yaml:"secrets"`

	Components SiteComponentConfigs `yaml:"components"`
}

func parseSitesNode(cfg *MachConfig, sitesNode *yaml.Node) error {
	if err := sitesNode.Decode(&cfg.Sites); err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}

	for _, site := range sitesNode.Content {
		nodes := MapYamlNodes(site.Content)
		siteId := nodes["identifier"].Value

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

			if err := plugin.SetSiteConfig(siteId, data); err != nil {
				return fmt.Errorf("%s.SetSiteConfig failed: %w", plugin.Name, err)
			}
		}

		if node, ok := nodes["endpoints"]; ok {
			cli.DeprecationWarning(&cli.DeprecationOptions{
				Message: "the usage of endpoints is deprecated and will be removed in the next major version",
				Details: `
				Configure endpoints by creating the API gateway as a separate component and configuring routes in the
				in the components that provide the logic.
				`,
				Site: siteId,
			})
			if err := parseSiteEndpointNode(cfg, siteId, node); err != nil {
				return fmt.Errorf("failed to parse endpoints: %w", err)
			}
		}

		if err := parseSiteComponentsNode(cfg, siteId, nodes["components"]); err != nil {
			return err
		}
	}

	for k, s := range cfg.Sites {
		if s.Deployment == nil {
			log.Debug().Msgf("No site deployment type specified for %s; defaulting to global setting", s.Identifier)
			var siteDeployment = cfg.MachComposer.Deployment
			cfg.Sites[k].Deployment = &siteDeployment
		}
	}

	return resolveSiteComponents(cfg)
}

func parseSiteEndpointNode(cfg *MachConfig, siteId string, node *yaml.Node) error {
	nodes := MapYamlNodes(node.Content)
	knownTags := []string{"url", "key", "zone", "throttling_rate_limit", "throttling_burst_limit"}

	for endpointId, endpointNode := range nodes {
		if endpointNode.Kind == yaml.ScalarNode {
			cli.DeprecationWarning(&cli.DeprecationOptions{
				Message: fmt.Sprintf("endpoint '%s' should be mapping n with a plugin tag", endpointId),
				Details: utils.TrimIndent(`
					For example instead of:
						endpoints:
							my-endpoint: my-url.example.org

					You should use:
						endpoints:
							my-endpoint:
								aws:
									url: my-url.example.org
				`),
			})

			cloudPlugin, err := cfg.Plugins.Get(cfg.Global.Cloud)
			if err != nil {
				return err
			}

			data := map[string]any{
				"url": endpointNode.Value,
			}

			if err := cloudPlugin.SetSiteEndpointConfig(siteId, endpointId, data); err != nil {
				return fmt.Errorf("cloudPlugin.SetSiteEndpointConfig: %w", err)
			}

			continue
		}

		children := MapYamlNodes(endpointNode.Content)

		if len(pie.Intersect(knownTags, pie.Keys(children))) > 0 {
			cli.DeprecationWarning(&cli.DeprecationOptions{
				Message: fmt.Sprintf("All endpoint properties on '%s' should be moved within the plugin tag", endpointId),
				Details: utils.TrimIndent(`
					For example instead of:
						endpoints:
							my-endpoint:
								url: my-url.example.org
								aws:
									throttling_burst_limit: 5000
									throttling_rate_limit: 10000

					You should use:
						endpoints:
							my-endpoint:
								aws:
									url: my-url.example.org
									throttling_burst_limit: 5000
									throttling_rate_limit: 10000
				`),
			})
		}

		legacyData := map[string]any{}
		for key, n := range children {
			if pie.Contains(knownTags, key) {
				legacyData[key] = n.Value
			}
		}

		for key, n := range children {
			if pie.Contains(knownTags, key) {
				continue
			}

			data, err := nodeAsMap(n)
			if err != nil {
				return err
			}

			// Copy the old legacy data in the resulting data. Can be removed
			// once we always require the nodes to be part of the plugin
			for k, v := range legacyData {
				if _, ok := data[k]; !ok {
					data[k] = v
				}
			}

			for _, plugin := range cfg.Plugins.All() {
				pluginNode, ok := nodes[plugin.Name]
				if ok {
					var err error
					data, err = nodeAsMap(pluginNode)
					if err != nil {
						return err
					}
				}

				if err := plugin.SetSiteEndpointConfig(siteId, endpointId, data); err != nil {
					return fmt.Errorf("%s.SetSiteEndpointConfig failed: %w", plugin.Name, err)
				}
			}
		}
	}
	return nil
}

func parseSiteComponentsNode(cfg *MachConfig, siteKey string, node *yaml.Node) error {
	// Exit early when no components are defined for this siteKey. Not a common
	// scenario, but still
	if node == nil {
		return nil
	}

	for _, component := range node.Content {
		nodes := MapYamlNodes(component.Content)
		componentKey := nodes["name"].Value

		migrateCommercetools(siteKey, componentKey, nodes)

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

			if err := plugin.SetSiteComponentConfig(siteKey, componentKey, data); err != nil {
				return err
			}
		}
	}
	return nil
}

func resolveSiteComponents(cfg *MachConfig) error {
	components := make(map[string]*ComponentConfig, len(cfg.Components))
	for i, c := range cfg.Components {
		components[c.Name] = &cfg.Components[i]
	}

	for _, site := range cfg.Sites {
		if len(site.Components) < 1 {
			continue
		}

		for i := range site.Components {
			c := &site.Components[i]

			if c.Deployment == nil {
				log.Debug().Msgf("No site component deployment type specified for %s; defaulting to global setting", c.Name)
				var siteComponentDeployment = *site.Deployment
				c.Deployment = &siteComponentDeployment
			}

			ref, ok := components[c.Name]
			if !ok {
				return fmt.Errorf("component %s does not exist in global components", c.Name)
			}
			c.Definition = ref
		}
	}
	return nil
}

// migrateCommercetools moves the store_variables and store_secrets under the
// commercetools node. Needed to say backwards compatible
func migrateCommercetools(site, name string, nodes map[string]*yaml.Node) {
	needsMigrate := false
	if _, ok := nodes["store_variables"]; ok {
		needsMigrate = true
	}
	if _, ok := nodes["store_secrets"]; ok {
		needsMigrate = true
	}
	if !needsMigrate {
		return
	}

	cli.DeprecationWarning(&cli.DeprecationOptions{
		Site:      site,
		Component: name,
		Message:   fmt.Sprintf("component %s in site %s is using deprecated syntax", name, site),
		Details: `
			The nodes 'store_variables' and 'store_secrets' are part of the
			commercetools plugin and should therefore be specified as such.

			For example instead of:
			 	variables:
			 		store_variables:
			 			MY_STORE:
			 				username: value
			 		store_secrets:
			 			MY_STORE:
			 				password: secret-value

			 You should use:
				commercetools:
					store_variables:
						MY_STORE:
							username: value
					store_secrets:
						MY_STORE:
							password: secret-value
		`,
		Version: "2.6",
	})

	if _, ok := nodes["commercetools"]; !ok {
		nodes["commercetools"] = &yaml.Node{
			Kind:    yaml.MappingNode,
			Tag:     "!!map",
			Content: []*yaml.Node{},
		}

		if val, ok := nodes["store_variables"]; ok {
			keyNode := &yaml.Node{
				Kind:  yaml.ScalarNode,
				Tag:   "!!str",
				Value: "store_variables",
			}
			nodes["commercetools"].Content = append(
				nodes["commercetools"].Content,
				keyNode,
				val,
			)
		}
		if val, ok := nodes["store_secrets"]; ok {
			keyNode := &yaml.Node{
				Kind:  yaml.ScalarNode,
				Tag:   "!!str",
				Value: "store_secrets",
			}
			nodes["commercetools"].Content = append(
				nodes["commercetools"].Content,
				keyNode,
				val,
			)
		}
	}
}
