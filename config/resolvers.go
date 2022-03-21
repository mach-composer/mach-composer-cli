package config

import (
	"log"
	"path/filepath"
	"strings"
)

func Process(cfg *Root) {

	// resolve_variables(config, config.variables, config.variables_encrypted)
	// parse_global_config(config)
	// resolve_component_definitions(config)
	ResolveComponentDefinitions(cfg)
	ResolveSiteConfigs(cfg)

}

func ResolveComponentDefinitions(cfg *Root) {
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
			if cfg.Global.Cloud == "aws" {
				c.Integrations = append(c.Integrations, "aws")
			} else if cfg.Global.Cloud == "azure" {
				c.Integrations = append(c.Integrations, "azure")
			}
		}
	}
}

func ResolveSiteConfigs(cfg *Root) {
	ResolveSentryConfig(cfg)
	ResolveSiteComponents(cfg)

	for i := range cfg.Sites {
		ResolveComponentEndpoints(&cfg.Sites[i])
	}
}

func ResolveSiteComponents(cfg *Root) {
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

func ResolveSentryConfig(cfg *Root) {
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
