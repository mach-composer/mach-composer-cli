package amplience

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/labd/mach-composer/internal/plugins/mcsdk"
	"github.com/labd/mach-composer/internal/plugins/shared"
)

func NewAmpliencePlugin() mcsdk.MachComposerPlugin {
	state := &AmpliencePlugin{
		provider:    "0.3.7",
		siteConfigs: map[string]*AmplienceConfig{},
	}

	return mcsdk.NewPlugin(&mcsdk.PluginSchema{
		Identifier: "amplience",

		Configure: state.Configure,
		IsEnabled: func() bool { return state.enabled },

		// Config
		SetGlobalConfig: state.SetGlobalConfig,
		SetSiteConfig:   state.SetSiteConfig,

		// Renders
		RenderTerraformProviders: state.TerraformRenderProviders,
		RenderTerraformResources: state.TerraformRenderResources,
		RenderTerraformComponent: state.RenderTerraformComponent,
	})
}

type AmpliencePlugin struct {
	environment  string
	provider     string
	globalConfig *AmplienceConfig
	siteConfigs  map[string]*AmplienceConfig
	enabled      bool
}

func (p *AmpliencePlugin) Configure(environment string, provider string) error {
	p.environment = environment
	if provider != "" {
		p.provider = provider
	}
	return nil
}

func (p *AmpliencePlugin) SetGlobalConfig(data map[string]any) error {
	cfg := AmplienceConfig{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	p.globalConfig = &cfg
	p.enabled = true
	return nil
}

func (p *AmpliencePlugin) SetSiteConfig(site string, data map[string]any) error {
	cfg := AmplienceConfig{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	p.siteConfigs[site] = &cfg
	p.enabled = true
	return nil
}

func (p *AmpliencePlugin) getSiteConfig(site string) *AmplienceConfig {
	result := &AmplienceConfig{}
	if p.globalConfig != nil {
		result.ClientID = p.globalConfig.ClientID
		result.ClientSecret = p.globalConfig.ClientSecret
		result.HubID = p.globalConfig.HubID
	}

	cfg, ok := p.siteConfigs[site]
	if ok {
		if cfg.ClientID != "" {
			result.ClientID = cfg.ClientID
		}
		if cfg.ClientSecret != "" {
			result.ClientSecret = cfg.ClientSecret
		}
		if cfg.HubID != "" {
			result.HubID = cfg.HubID
		}
	}

	if result.ClientID == "" {
		return nil
	}
	return result
}

func (p *AmpliencePlugin) TerraformRenderProviders(site string) (string, error) {
	result := fmt.Sprintf(`
	amplience = {
		source = "labd/amplience"
		version = "%s"
	}`, shared.VersionConstraint(p.provider))
	return result, nil
}

func (p *AmpliencePlugin) TerraformRenderResources(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	template := `
		provider "amplience" {
			client_id        = {{ .ClientID|printf "%q" }}
			client_secret    = {{ .ClientSecret|printf "%q" }}
			hub_id           = {{ .HubID|printf "%q" }}
		}
	`
	return shared.RenderGoTemplate(template, cfg)
}

func (p *AmpliencePlugin) RenderTerraformComponent(site string, component string) (*mcsdk.ComponentSnippets, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return nil, nil
	}

	template := `
		amplience_client_id = {{ .ClientID|printf "%q" }}
		amplience_client_secret = {{ .ClientSecret|printf "%q" }}
		amplience_hub_id = {{ .HubID|printf "%q" }}
	`
	vars, err := shared.RenderGoTemplate(template, cfg)
	if err != nil {
		return nil, err
	}
	result := &mcsdk.ComponentSnippets{
		Variables: vars,
	}
	return result, nil
}
