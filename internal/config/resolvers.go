package config

import (
	"log"
	"path/filepath"
	"strings"
)

func ProcessConfig(cfg *MachConfig) error {
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
		c.Integrations = append(c.Integrations, cfg.Global.Cloud)
	}
	return c, nil
}

func resolveSiteConfigs(cfg *MachConfig) error {
	resolveSiteComponents(cfg)
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
		}
	}
}
