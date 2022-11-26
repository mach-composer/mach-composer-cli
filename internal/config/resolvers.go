package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ProcessConfig(cfg *MachConfig) error {
	// resolve_variables(config, config.variables, config.variables_encrypted)
	// parse_global_config(config)
	// resolve_component_definitions(config)
	if err := resolveComponentDefinitions(cfg); err != nil {
		return err
	}
	if err := resolveSiteConfigs(cfg); err != nil {
		return err
	}
	return nil
}

func resolveComponentDefinitions(cfg *MachConfig) error {
	for i := range cfg.Components {
		if _, err := resolveComponentDefinition(&cfg.Components[i], cfg); err != nil {
			return err
		}
	}
	return nil
}

func resolveComponentDefinition(c *Component, cfg *MachConfig) (*Component, error) {
	// Terraform needs absolute paths to modules
	if strings.HasPrefix(c.Source, ".") {
		if val, err := filepath.Abs(c.Source); err == nil {
			c.Source = val
		} else {
			return nil, err
		}
	}

	// If no integrations are given, set the Cloud integrations as default
	if len(c.Integrations) < 1 {
		if cfg.Global.Cloud == AWS {
			c.Integrations = append(c.Integrations, AWS)
		} else if cfg.Global.Cloud == Azure {
			c.Integrations = append(c.Integrations, Azure)
		}
	}

	if cfg.Global.Cloud == Azure {
		c.Azure = &ComponentAzureConfig{}
	}

	if c.Azure != nil && c.Azure.ShortName == "" {
		c.Azure.ShortName = c.Name
	}

	return c, nil
}

func resolveSiteConfigs(cfg *MachConfig) error {
	resolveAzureConfig(cfg)
	resolveSentryConfig(cfg)
	resolveSiteComponents(cfg)

	for i := range cfg.Sites {
		err := resolveComponentEndpoints(&cfg.Sites[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func resolveSiteComponents(cfg *MachConfig) {
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
				log.Fatalf("Component %s does not exist in global components.", c.Name)
			}
			c.Definition = ref

			if site.Sentry != nil {
				if c.Sentry == nil {
					c.Sentry = newSentryConfig(site.Sentry)
				} else {
					c.Sentry.merge(site.Sentry)
				}
			}
		}
	}
}

func resolveSentryConfig(cfg *MachConfig) {
	if cfg.Global.SentryConfig != nil {
		for i := range cfg.Sites {
			s := &cfg.Sites[i]
			if s.Sentry == nil {
				s.Sentry = newSentryConfigFromGlobal(cfg.Global.SentryConfig)
			} else {
				s.Sentry.mergeGlobal(cfg.Global.SentryConfig)
			}
		}
	}
}

func resolveAzureConfig(cfg *MachConfig) {
	if cfg.Global.Cloud != "azure" {
		return
	}

	if cfg.Global.SentryConfig != nil {
		for i := range cfg.Sites {
			s := &cfg.Sites[i]

			if s.Azure == nil {
				s.Azure = &SiteAzureSettings{}
			}
			s.Azure.merge(cfg.Global.Azure)
			if s.Azure.ResourceGroup != "" {
				fmt.Fprintf(
					os.Stderr,
					"WARNING: resource_group on %s is used (%s). "+
						"Make sure it wasn't managed by MACH before otherwise "+
						"the resource group will get deleted.",
					s.Identifier, s.Azure.ResourceGroup,
				)
			}
		}
	}
}

func resolveComponentEndpoints(site *Site) error {
	if err := site.resolveEndpoints(); err != nil {
		return err
	}

	components := site.EndpointComponents()
	for i := range site.Endpoints {
		ep := &site.Endpoints[i]
		if c, ok := components[ep.Key]; ok {
			ep.Components = c
		} else {
			ep.Components = make([]SiteComponent, 0)
		}
	}
	return nil
}

func stringContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
