package sentry

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/labd/mach-composer/internal/plugins/shared"
)

type SentryPlugin struct {
	environment  string
	provider     string
	globalConfig GlobalConfig
	siteConfigs  map[string]*SiteConfig
}

func NewSentryPlugin() *SentryPlugin {
	return &SentryPlugin{
		siteConfigs: map[string]*SiteConfig{},
	}
}

func (p *SentryPlugin) Configure(environment string, provider string) error {
	p.environment = environment
	if provider != "" {
		p.provider = provider
	}
	return nil
}

func (p *SentryPlugin) IsEnabled() bool {
	return p.globalConfig.AuthToken != ""
}

func (p *SentryPlugin) Identifier() string {
	return "sentry"
}

func (p *SentryPlugin) SetRemoteStateBackend(data map[string]any) error {
	return fmt.Errorf("not supported by this plugin")
}

func (p *SentryPlugin) SetGlobalConfig(data map[string]any) error {
	if err := mapstructure.Decode(data, &p.globalConfig); err != nil {
		return err
	}
	return nil
}

func (p *SentryPlugin) SetSiteConfig(site string, data map[string]any) error {
	cfg := SiteConfig{
		Components: map[string]ComponentConfig{},
	}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	p.siteConfigs[site] = &cfg
	return nil
}

func (p *SentryPlugin) SetSiteComponentConfig(site string, component string, data map[string]any) error {
	cfg := ComponentConfig{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}

	siteCfg, ok := p.siteConfigs[site]
	if !ok {
		siteCfg = &SiteConfig{
			Components: map[string]ComponentConfig{},
		}
		p.siteConfigs[site] = siteCfg
	}
	siteCfg.Components[component] = cfg

	return nil
}

func (p *SentryPlugin) SetSiteEndpointsConfig(site string, data map[string]any) error {
	return nil
}

func (p *SentryPlugin) SetComponentConfig(component string, data map[string]any) error {
	return nil
}

func (p *SentryPlugin) SetComponentEndpointsConfig(component string, endpoints map[string]string) error {
	return nil
}

func (p *SentryPlugin) TerraformRenderStateBackend(site string) (string, error) {
	return "", nil
}

func (p *SentryPlugin) TerraformRenderProviders(site string) (string, error) {
	result := fmt.Sprintf(`
		sentry = {
			source = "jianyuan/sentry"
			version = "%s"
		}`, p.provider)
	return result, nil
}

func (p *SentryPlugin) TerraformRenderResources(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	template := `
		provider "sentry" {
			token = {{ .AuthToken|printf "%q" }}
			base_url = {{ if .BaseURL }}{{ .BaseURL|printf "%q" }}{{ else }}"https://sentry.io/api/"{{ end }}
		}
	`
	return shared.RenderGoTemplate(template, p.globalConfig)
}

func (p *SentryPlugin) TerraformRenderComponentResources(site string, component string) (string, error) {
	if p.globalConfig.AuthToken == "" {
		return "", nil
	}

	cfg := p.getComponentSiteConfig(site, component)
	if cfg == nil {
		return "", nil
	}

	templateContext := struct {
		ComponentName string
		SiteName      string
		Global        *GlobalConfig
		Config        *ComponentConfig
	}{
		ComponentName: component,
		SiteName:      site,
		Global:        &p.globalConfig,
		Config:        cfg,
	}

	template := `
		resource "sentry_key" "{{ .ComponentName }}" {
		organization      = {{ .Global.Organization|printf "%q" }}
		project           = {{ .Config.Project|printf "%q" }}
		name              = "{{ .SiteName }}_{{ .ComponentName }}"
		{{ if .Config.RateLimitWindow }}
		rate_limit_window = {{ .Config.RateLimitWindow }}
		{{ end }}
		{{ if .Config.RateLimitCount }}
		rate_limit_count  = {{ .Config.RateLimitCount }}
		{{ end }}
		}
	`
	return shared.RenderGoTemplate(template, templateContext)
}

func (p *SentryPlugin) TerraformRenderComponentVars(site string, component string) (string, error) {
	if p.globalConfig.AuthToken == "" {
		return "", nil
	}

	cfg := p.getComponentSiteConfig(site, component)
	if cfg == nil {
		return "", nil
	}

	templateContext := struct {
		ComponentName string
		SiteName      string
		Global        *GlobalConfig
		Config        *ComponentConfig
	}{
		ComponentName: component,
		SiteName:      site,
		Global:        &p.globalConfig,
		Config:        cfg,
	}

	template := `
		sentry_dsn = {{ if .Global.AuthToken }}sentry_key.{{ .ComponentName }}.dsn_secret{{ else }}"{{ .Config.DSN }}"{{ end }}
	`
	return shared.RenderGoTemplate(template, templateContext)
}

func (p *SentryPlugin) TerraformRenderComponentDependsOn(site string, component string) ([]string, error) {
	return []string{}, nil // TODO. sentry_key.component
}

func (p *SentryPlugin) TerraformRenderComponentProviders(site string, component string) ([]string, error) {
	return []string{}, nil
}

func (p *SentryPlugin) getSiteConfig(site string) *SiteConfig {
	cfg, ok := p.siteConfigs[site]
	if !ok {
		return &SiteConfig{}
	}
	return cfg.extendGlobalConfig(&p.globalConfig)
}

func (p *SentryPlugin) getComponentSiteConfig(site, name string) *ComponentConfig {
	siteCfg := p.getSiteConfig(site)
	if siteCfg == nil {
		return nil
	}
	return siteCfg.getComponentSiteConfig(name)
}
