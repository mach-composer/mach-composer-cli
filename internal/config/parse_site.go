package config

import (
	"fmt"

	"github.com/elliotchance/pie/v2"
	"gopkg.in/yaml.v3"

	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/mach-composer/mach-composer-cli/internal/variables"
)

func parseSitesNode(cfg *MachConfig, sitesNode *yaml.Node) error {
	if err := sitesNode.Decode(&cfg.Sites); err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}

	for _, site := range sitesNode.Content {
		nodes := mapYamlNodes(site.Content)
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
			if err := parseSiteEndpointNode(cfg, siteId, node); err != nil {
				return fmt.Errorf("failed to parse endpoints: %w", err)
			}
		}

		if err := parseSiteComponentsNode(cfg, siteId, nodes["components"]); err != nil {
			return err
		}
	}

	return resolveSiteComponents(cfg)
}

func parseSiteEndpointNode(cfg *MachConfig, siteId string, node *yaml.Node) error {
	nodes := mapYamlNodes(node.Content)
	knownTags := []string{"url", "key", "zone", "throttling_rate_limit", "throttling_burst_limit"}

	for endpointId, endpointNode := range nodes {
		if endpointNode.Kind == yaml.ScalarNode {
			cli.DeprecationWarning(&cli.DeprecationOptions{
				Message: fmt.Sprintf("endpoint '%s' should be mapping node with a plugin tag", endpointId),
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

		childs := mapYamlNodes(endpointNode.Content)

		if len(pie.Intersect(knownTags, pie.Keys(childs))) > 0 {
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
		for key, node := range childs {
			if pie.Contains(knownTags, key) {
				legacyData[key] = node.Value
			}
		}

		for key, node := range childs {
			if pie.Contains(knownTags, key) {
				continue
			}

			data, err := nodeAsMap(node)
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

			if err := cfg.Plugins.SetSiteEndpointConfig(key, siteId, endpointId, data); err != nil {
				return fmt.Errorf("plugins.SetSiteEndpointConfig: %w", err)
			}
		}
	}
	return nil
}

func parseSiteComponentsNode(cfg *MachConfig, site string, node *yaml.Node) error {
	// Exit early when no components are defined for this site. Not a common
	// scenario, but still
	if node == nil {
		return nil
	}

	for _, component := range node.Content {
		nodes := mapYamlNodes(component.Content)
		identifier := nodes["name"].Value

		migrateCommercetools(site, identifier, nodes)

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

			if err := plugin.SetSiteComponentConfig(site, identifier, data); err != nil {
				return err
			}
		}
	}
	return nil
}

func resolveSiteComponents(cfg *MachConfig) error {
	components := make(map[string]*Component, len(cfg.Components))
	for i, c := range cfg.Components {
		components[c.Name] = &cfg.Components[i]
	}

	for _, site := range cfg.Sites {
		if len(site.Components) < 1 {
			continue
		}

		for i := range site.Components {
			c := &site.Components[i]

			ref, ok := components[c.Name]
			if !ok {
				return fmt.Errorf("component %s does not exist in global components.", c.Name)
			}
			c.Definition = ref

			var err error
			c.Variables, err = variables.InterpolateComponentVars(c.Variables)
			if err != nil {
				return err
			}

			c.Secrets, err = variables.InterpolateComponentVars(c.Secrets)
			if err != nil {
				return err
			}
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
