package amplience

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/mitchellh/mapstructure"
)

type AmpliencePlugin struct {
	globalConfig *AmplienceConfig
	siteConfigs  map[string]*AmplienceConfig
	enabled      bool
}

func NewAmpliencePlugin() *AmpliencePlugin {
	return &AmpliencePlugin{
		siteConfigs: map[string]*AmplienceConfig{},
	}
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

func (p *AmpliencePlugin) TerraformRenderStateBackend(site string) string {
	return ""
}

func (p *AmpliencePlugin) TerraformRenderProviders(site string) string {
	return `
	amplience = {
		source = "labd/amplience"
		version = "0.3.7"
	}`
}

func (p *AmpliencePlugin) TerraformRenderResources(site string) string {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return ""
	}

	template := `
		provider "amplience" {
			client_id        = {{ .ClientID|printf "%q" }}
			client_secret    = {{ .ClientSecret|printf "%q" }}
			hub_id           = {{ .HubID|printf "%q" }}
		}
	`
	return renderTemplate(template, cfg)
}

func (p *AmpliencePlugin) TerraformRenderComponentResources(site string, component string) string {
	return ""
}

func (p *AmpliencePlugin) TerraformRenderComponentVars(site, component string) string {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return ""
	}

	template := `
		amplience_client_id = {{ .ClientID|printf "%q" }}
		amplience_client_secret = {{ .ClientSecret|printf "%q" }}
		amplience_hub_id = {{ .HubID|printf "%q" }}
	`
	return renderTemplate(template, cfg)
}

func (p *AmpliencePlugin) TerraformRenderComponentDependsOn(site string, component string) []string {
	return []string{}
}

func (p *AmpliencePlugin) TerraformRenderComponentProviders(site string, component string) []string {
	return []string{}
}

func renderTemplate(t string, data any) string {
	tpl, err := template.New("template").Parse(t)
	if err != nil {
		panic(err)
	}

	var content bytes.Buffer
	if err := tpl.Execute(&content, data); err != nil {
		panic(err)
	}
	return content.String()
}
