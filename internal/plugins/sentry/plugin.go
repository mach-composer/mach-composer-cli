package sentry

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/labd/mach-composer/internal/plugins/mcsdk"
	"github.com/labd/mach-composer/internal/plugins/shared"
)

type SentryPlugin struct {
	environment  string
	provider     string
	globalConfig GlobalConfig
	siteConfigs  map[string]*SiteConfig
}

func NewSentryPlugin() mcsdk.MachComposerPlugin {
	state := &SentryPlugin{
		provider:    "0.6.0",
		siteConfigs: map[string]*SiteConfig{},
	}

	return mcsdk.NewPlugin(&mcsdk.PluginSchema{
		Identifier: "sentry",

		Configure: state.Configure,
		IsEnabled: state.IsEnabled,

		// Config
		SetGlobalConfig:        state.SetGlobalConfig,
		SetSiteConfig:          state.SetSiteConfig,
		SetSiteComponentConfig: state.SetSiteComponentConfig,

		// Renders
		RenderTerraformProviders: state.TerraformRenderProviders,
		RenderTerraformResources: state.TerraformRenderResources,
		RenderTerraformComponent: state.RenderTerraformComponent,
	})
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

func (p *SentryPlugin) TerraformRenderProviders(site string) (string, error) {
	result := fmt.Sprintf(`
		sentry = {
			source = "jianyuan/sentry"
			version = "%s"
		}`, shared.VersionConstraint(p.provider))
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

func (p *SentryPlugin) RenderTerraformComponent(site string, component string) (*mcsdk.ComponentSnippets, error) {
	cfg := p.getComponentSiteConfig(site, component)

	vars := fmt.Sprintf("sentry_dsn = \"%s\"", cfg.DSN)
	if p.globalConfig.AuthToken != "" {
		vars = fmt.Sprintf("sentry_dsn = sentry_key.%s.dsn_secret", component)
	}

	resources, err := terraformRenderComponentResources(site, component, cfg, &p.globalConfig)
	if err != nil {
		return nil, err
	}

	result := &mcsdk.ComponentSnippets{
		Variables: vars,
		Resources: resources,
	}
	return result, nil
}

func (p *SentryPlugin) getSiteConfig(site string) *SiteConfig {
	cfg, ok := p.siteConfigs[site]
	if !ok {
		cfg = &SiteConfig{}
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

func terraformRenderComponentResources(site, component string, cfg *ComponentConfig, globalCfg *GlobalConfig) (string, error) {
	templateContext := struct {
		ComponentName string
		SiteName      string
		Global        *GlobalConfig
		Config        *ComponentConfig
	}{
		ComponentName: component,
		SiteName:      site,
		Global:        globalCfg,
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
