package commercetools

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/creasty/defaults"
	"github.com/mitchellh/mapstructure"
)

type CommercetoolsPlugin struct {
	siteConfigs map[string]*SiteConfig
}

func NewCommercetoolsPlugin() *CommercetoolsPlugin {
	return &CommercetoolsPlugin{
		siteConfigs: map[string]*SiteConfig{},
	}
}

func (p *CommercetoolsPlugin) IsEnabled() bool {
	return len(p.siteConfigs) > 0
}

func (p *CommercetoolsPlugin) Identifier() string {
	return "sentry"
}

func (p *CommercetoolsPlugin) SetRemoteStateBackend(data map[string]any) error {
	return fmt.Errorf("not supported by this plugin")
}

func (p *CommercetoolsPlugin) SetGlobalConfig(data map[string]any) error {
	return nil
}

func (p *CommercetoolsPlugin) SetSiteConfig(site string, data map[string]any) error {
	cfg := SiteConfig{
		Components: map[string]ComponentConfig{},
	}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}

	if cfg.Frontend != nil {
		fmt.Println("[WARN] Site", site, "commercetools frontend block is deprecated and will be removed soon")
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

func (p *CommercetoolsPlugin) SetSiteEndpointsConfig(site string, data map[string]any) error {
	return nil
}

func (p *CommercetoolsPlugin) SetComponentConfig(component string, data map[string]any) error {
	return nil
}

func (p *CommercetoolsPlugin) SetComponentEndpointsConfig(component string, endpoints map[string]string) error {
	return nil
}

func (p *CommercetoolsPlugin) getSiteConfig(site string) *SiteConfig {
	cfg, ok := p.siteConfigs[site]
	if !ok {
		return nil
	}
	return cfg
}

func (p *CommercetoolsPlugin) TerraformRenderStateBackend(site string) string {
	return ""
}

func (p *CommercetoolsPlugin) TerraformRenderProviders(site string) string {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return ""
	}

	return `
    commercetools = {
      source = "labd/commercetools"
      version = "0.30.0"
    }
	`
}

func (p *CommercetoolsPlugin) TerraformRenderResources(site string) string {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return ""
	}

	content, err := Render(cfg)
	if err != nil {
		panic(err)
	}

	return content
}

func (p *CommercetoolsPlugin) TerraformRenderComponentResources(site string, component string) string {
	return ""
}

func (p *CommercetoolsPlugin) TerraformRenderComponentVars(site string, component string) string {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return ""
	}

	templateContext := struct {
		Site      *SiteConfig
		Component *ComponentConfig
	}{
		Site:      cfg,
		Component: cfg.getComponentSiteConfig(component),
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
	return renderTemplate(template, templateContext)
}

func (p *CommercetoolsPlugin) TerraformRenderComponentDependsOn(site string, component string) []string {
	return []string{"null_resource.commercetools"}
}

func renderTemplate(t string, data any) string {
	tpl, err := template.New("template-1").Parse(t)
	if err != nil {
		panic(err)
	}

	var content bytes.Buffer
	if err := tpl.Execute(&content, data); err != nil {
		panic(err)
	}
	return content.String()
}
