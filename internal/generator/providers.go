package generator

import (
	"embed"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

//go:embed templates/provider.tmpl
var providerTmpl embed.FS

type RenderProviderConfiguration struct {
	Name          string
	Configuration map[string]any `yaml:"configuration,omitempty"`
}

func renderProviders(cfg *config.MachConfig) ([]string, error) {
	tpl, err := providerTmpl.ReadFile("templates/provider.tmpl")
	if err != nil {
		return nil, err
	}

	siteConfigs := cfg.Global.TerraformConfig.ProviderConfigs

	var renderedProviders []string
	for _, provider := range siteConfigs {
		data := RenderProviderConfiguration{
			Name:          provider.Name,
			Configuration: provider.Configuration,
		}

		renderedProvider, err := utils.RenderGoTemplate(string(tpl), data)
		if err != nil {
			return nil, err
		}

		renderedProviders = append(renderedProviders, renderedProvider)
	}

	return renderedProviders, nil
}
