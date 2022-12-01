package commercetools

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"

	"github.com/labd/mach-composer/internal/plugins/mcsdk"
	"github.com/labd/mach-composer/internal/plugins/shared"
)

type CommercetoolsPlugin struct {
	environment string
	provider    string
	siteConfigs map[string]*SiteConfig
}

func NewCommercetoolsPlugin() mcsdk.MachComposerPlugin {
	state := &CommercetoolsPlugin{
		provider:    "0.30.0",
		siteConfigs: map[string]*SiteConfig{},
	}

	return mcsdk.NewPlugin(&mcsdk.PluginSchema{
		Identifier: "commercetools",

		Configure: state.Configure,
		IsEnabled: state.IsEnabled,

		// Config
		SetSiteConfig: state.SetSiteConfig,

		// Renders
		RenderTerraformProviders: state.TerraformRenderProviders,
		RenderTerraformResources: state.TerraformRenderResources,
		RenderTerraformComponent: state.RenderTerraformComponent,
	})
}

func (p *CommercetoolsPlugin) Configure(environment string, provider string) error {
	p.environment = environment
	if provider != "" {
		p.provider = provider
	}
	return nil
}

func (p *CommercetoolsPlugin) IsEnabled() bool {
	return len(p.siteConfigs) > 0
}

func (p *CommercetoolsPlugin) SetSiteConfig(site string, data map[string]any) error {
	cfg := SiteConfig{
		Components: map[string]ComponentConfig{},
	}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}

	if cfg.Frontend != nil {
		logrus.Warnf("%s: commercetools frontend block is deprecated and will be removed soon\n", site)
	}

	if err := defaults.Set(&cfg); err != nil {
		return err
	}
	p.siteConfigs[site] = &cfg

	return nil
}

func (p *CommercetoolsPlugin) SetSiteComponentConfig(site string, component string, data map[string]any) error {
	siteConfig := p.getSiteConfig(site)
	if siteConfig == nil {
		return fmt.Errorf("no site config found")
	}

	cfg := ComponentConfig{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	siteConfig.Components[component] = cfg

	return nil
}

func (p *CommercetoolsPlugin) TerraformRenderProviders(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	result := fmt.Sprintf(`
		commercetools = {
		source = "labd/commercetools"
		version = "%s"
		}
	`, shared.VersionConstraint(p.provider))
	return result, nil
}

func (p *CommercetoolsPlugin) TerraformRenderResources(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	content, err := renderResources(cfg)
	if err != nil {
		return "", err
	}

	return content, nil
}

func (p *CommercetoolsPlugin) RenderTerraformComponent(site string, component string) (*mcsdk.ComponentSnippets, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return nil, nil
	}
	componentCfg := cfg.getComponentSiteConfig(component)

	vars, err := terraformRenderComponentVars(cfg, componentCfg)
	if err != nil {
		return nil, err
	}

	result := &mcsdk.ComponentSnippets{
		Variables: vars,
		DependsOn: []string{"null_resource.commercetools"},
	}
	return result, nil
}

func (p *CommercetoolsPlugin) getSiteConfig(site string) *SiteConfig {
	cfg, ok := p.siteConfigs[site]
	if !ok {
		return nil
	}
	return cfg
}

func terraformRenderComponentVars(cfg *SiteConfig, componentCfg *ComponentConfig) (string, error) {
	templateContext := struct {
		Site      *SiteConfig
		Component *ComponentConfig
	}{
		Site:      cfg,
		Component: componentCfg,
	}

	template := `
		ct_project_key    = {{ .Site.ProjectKey|printf "%q" }}
		ct_api_url        = {{ .Site.APIURL|printf "%q" }}
		ct_auth_url       = {{ .Site.TokenURL|printf "%q" }}

		ct_stores = {
			{{ $component := .Component }}
			{{ range $store := .Site.Stores }}
				{{ $store.Key }} =  {
					key = {{ $store.Key|printf "%q" }}
					variables = {
						{{ range $key, $value := index $component.StoreVariables $store.Key }}
						{{ $key }} = {{ $value|printf "%#v" }}
						{{ end }}
					}
					secrets = {
						{{ range $key, $value := index $component.StoreSecrets $store.Key }}
						{{ $key }} = {{ $value|printf "%#v" }}
						{{ end }}
					}
				}
			{{ end }}
		}
	`
	return shared.RenderGoTemplate(template, templateContext)
}
