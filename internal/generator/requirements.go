package generator

import (
	"embed"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-plugin-sdk/v2/helpers"
)

//go:embed templates/requirement.tmpl
var requirementsTmpl embed.FS

type RenderProviderRequirements struct {
	Name    string
	Source  string
	Version string
}

func renderRequirements(cfg *config.MachConfig) ([]string, error) {
	tpl, err := requirementsTmpl.ReadFile("templates/requirement.tmpl")
	if err != nil {
		return nil, err
	}

	siteConfigs := cfg.Global.TerraformConfig.ProviderConfigs

	var renderedRequirements []string
	for _, provider := range siteConfigs {
		providerData := RenderProviderRequirements{
			Name:    provider.Name,
			Source:  provider.Source,
			Version: provider.Version,
		}
		renderedRequirement, err := helpers.RenderGoTemplate(string(tpl), providerData)
		if err != nil {
			return nil, err
		}

		renderedRequirements = append(renderedRequirements, renderedRequirement)
	}

	return renderedRequirements, nil
}
