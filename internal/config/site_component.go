package config

import (
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/mach-composer/mach-composer-cli/internal/config/variable"
	"github.com/rs/zerolog/log"
)

type SiteComponentConfigs []SiteComponentConfig

func (s *SiteComponentConfigs) Get(name string) (*SiteComponentConfig, error) {
	for _, site := range *s {
		if site.Name == name {
			return &site, nil
		}
	}
	return nil, fmt.Errorf("site component %s not found", name)
}

type SiteComponentConfig struct {
	Name       string                `yaml:"name"`
	Definition *ComponentConfig      `yaml:"-"`
	Variables  variable.VariablesMap `yaml:"variables"`
	Secrets    variable.VariablesMap `yaml:"secrets"`
	Deployment *Deployment           `yaml:"deployment"`
	Count      string                `yaml:"count"`

	DependsOn []string `yaml:"depends_on"`
}

func (sc *SiteComponentConfig) HasCloudIntegration(g *GlobalConfig) bool {
	if sc.Definition == nil {
		log.Fatal().Msgf("ComponentConfig %s was not resolved properly (missing definition)", sc.Name)
	}
	return pie.Contains(sc.Definition.Integrations, g.Cloud)
}
