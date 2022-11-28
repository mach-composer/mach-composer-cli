package amplience

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/labd/mach-composer/internal/plugins/shared"
)

type AmpliencePlugin struct {
	environment  string
	globalConfig *AmplienceConfig
	siteConfigs  map[string]*AmplienceConfig
	enabled      bool
}

func NewAmpliencePlugin() *AmpliencePlugin {
	return &AmpliencePlugin{
		siteConfigs: map[string]*AmplienceConfig{},
	}
}

func (p *AmpliencePlugin) Initialize(environment string) error {
	p.environment = environment
	return nil
}

func (p *AmpliencePlugin) IsEnabled() bool {
	return p.enabled
}

func (p *AmpliencePlugin) Identifier() string {
	return "amplience"
}

func (p *AmpliencePlugin) SetRemoteStateBackend(data map[string]any) error {
	return fmt.Errorf("not supported by this plugin")
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

func (p *AmpliencePlugin) SetSiteComponentConfig(site string, component string, data map[string]any) error {
	return nil
}

func (p *AmpliencePlugin) SetSiteEndpointsConfig(site string, data map[string]any) error {
	return nil
}

func (p *AmpliencePlugin) SetComponentConfig(component string, data map[string]any) error {
	return nil
}

func (p *AmpliencePlugin) SetComponentEndpointsConfig(component string, endpoints map[string]string) error {
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

func (p *AmpliencePlugin) TerraformRenderStateBackend(site string) (string, error) {
	return "", nil
}

func (p *AmpliencePlugin) TerraformRenderProviders(site string) (string, error) {
	return `
	amplience = {
		source = "labd/amplience"
		version = "0.3.7"
	}`, nil
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

func (p *AmpliencePlugin) TerraformRenderComponentResources(site string, component string) (string, error) {
	return "", nil
}

func (p *AmpliencePlugin) TerraformRenderComponentVars(site, component string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	template := `
		amplience_client_id = {{ .ClientID|printf "%q" }}
		amplience_client_secret = {{ .ClientSecret|printf "%q" }}
		amplience_hub_id = {{ .HubID|printf "%q" }}
	`
	return shared.RenderGoTemplate(template, cfg)
}

func (p *AmpliencePlugin) TerraformRenderComponentDependsOn(site string, component string) ([]string, error) {
	return []string{}, nil
}

func (p *AmpliencePlugin) TerraformRenderComponentProviders(site string, component string) ([]string, error) {
	return []string{}, nil
}
