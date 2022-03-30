package config

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func Process(cfg *MachConfig) {

	// resolve_variables(config, config.variables, config.variables_encrypted)
	// parse_global_config(config)
	// resolve_component_definitions(config)
	ResolveComponentDefinitions(cfg)
	ResolveSiteConfigs(cfg)

}

func ResolveComponentDefinitions(cfg *MachConfig) {
	for _, c := range cfg.Components {

		// Terraform needs absolute paths to modules
		if strings.HasPrefix(c.Source, ".") {
			if val, err := filepath.Abs(c.Source); err == nil {
				c.Source = val
			} else {
				panic(err)
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

		if c.Azure.ShortName == "" {
			c.Azure.ShortName = c.Name
		}

	}
}

func ResolveSiteConfigs(cfg *MachConfig) {
	ResolveAzureConfig(cfg)
	ResolveSentryConfig(cfg)
	ResolveSiteComponents(cfg)

	for i := range cfg.Sites {
		ResolveComponentEndpoints(&cfg.Sites[i])
	}
}

func ResolveSiteComponents(cfg *MachConfig) {
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
					c.Sentry = NewSentryConfig(site.Sentry)
				} else {
					c.Sentry.Merge(site.Sentry)

				}
			}
		}
	}
}

func ResolveSentryConfig(cfg *MachConfig) {
	if cfg.Global.SentryConfig != nil {
		for i := range cfg.Sites {
			s := &cfg.Sites[i]
			if s.Sentry == nil {
				s.Sentry = NewSentryConfigFromGlobal(cfg.Global.SentryConfig)
			} else {
				s.Sentry.MergeGlobal(cfg.Global.SentryConfig)
			}
		}
	}
}

func ResolveAzureConfig(cfg *MachConfig) {
	if cfg.Global.Cloud != "azure" {
		return
	}

	if cfg.Global.SentryConfig != nil {
		for i := range cfg.Sites {
			s := &cfg.Sites[i]

			if s.Azure == nil {
				s.Azure = &SiteAzureSettings{}
			}
			s.Azure.Merge(cfg.Global.Azure)
			if s.Azure.ResourceGroup != "" {
				logrus.Errorf(
					"WARNING: resource_group on %s is used (%s). "+
						"Make sure it wasn't managed by MACH before otherwise "+
						"the resource group will get deleted.",
					s.Identifier, s.Azure.ResourceGroup,
				)
			}
		}
	}
}

func ResolveComponentEndpoints(site *Site) {
	site.ResolveEndpoints()

	components := site.EndpointComponents()
	for i := range site.Endpoints {
		ep := &site.Endpoints[i]
		if c, ok := components[ep.Key]; ok {
			ep.Components = c
		} else {
			ep.Components = make([]SiteComponent, 0)
		}
	}
}

func stringContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
