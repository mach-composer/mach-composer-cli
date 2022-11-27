package contentful

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/mitchellh/mapstructure"
)

type ContentfulPlugin struct {
	globalConfig *ContentfulConfig
	siteConfigs  map[string]*ContentfulConfig
	enabled      bool
}

func NewContentfulPlugin() *ContentfulPlugin {
	return &ContentfulPlugin{
		siteConfigs: map[string]*ContentfulConfig{},
	}
}

func (p *ContentfulPlugin) IsEnabled() bool {
	return p.enabled
}

func (p *ContentfulPlugin) Identifier() string {
	return "contentful"
}

func (p *ContentfulPlugin) SetRemoteStateBackend(data map[string]any) error {
	return fmt.Errorf("not supported by this plugin")
}

func (p *ContentfulPlugin) SetGlobalConfig(data map[string]any) error {
	cfg := ContentfulConfig{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	p.globalConfig = &cfg
	p.enabled = true
	return nil
}

func (p *ContentfulPlugin) SetSiteConfig(site string, data map[string]any) error {
	cfg := ContentfulConfig{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	p.siteConfigs[site] = &cfg
	p.enabled = true
	return nil
}

func (p *ContentfulPlugin) SetSiteComponentConfig(site string, component string, data map[string]any) error {
	return nil
}

func (p *ContentfulPlugin) SetSiteEndpointsConfig(site string, data map[string]any) error {
	return nil
}

func (p *ContentfulPlugin) SetComponentConfig(component string, data map[string]any) error {
	return nil
}

func (p *ContentfulPlugin) SetComponentEndpointsConfig(component string, endpoints map[string]string) error {
	return nil
}

func (p *ContentfulPlugin) TerraformRenderStateBackend(site string) string {
	return ""
}

func (p *ContentfulPlugin) TerraformRenderProviders(site string) string {
	return `
	contentful = {
		source = "labd/contentful"
		version = "0.1.0"
	}`
}

func (p *ContentfulPlugin) TerraformRenderResources(site string) string {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return ""
	}

	template := `
		provider "contentful" {
			cma_token       = {{ .CMAToken|printf "%q" }}
			organization_id = {{ .OrganizationID|printf "%q" }}
		  }

		  resource "contentful_space" "space" {
			name           = {{ .Space|printf "%q" }}
			default_locale = {{ .DefaultLocale|printf "%q" }}
		  }

		  resource "contentful_apikey" "apikey" {
			space_id = contentful_space.space.id

			name        = "frontend"
			description = "MACH generated frontend API key"
		  }

		  output "contentful_space_id" {
			value = contentful_space.space.id
		  }

		  output "contentful_apikey_access_token" {
			value = contentful_apikey.apikey.access_token
		  }
	`
	return renderTemplate(template, cfg)
}

func (p *ContentfulPlugin) TerraformRenderComponentResources(site string, component string) string {
	return ""
}

func (p *ContentfulPlugin) TerraformRenderComponentVars(site, component string) string {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return ""
	}

	return `
    	contentful_space_id = contentful_space.space.id
	`
}

func (p *ContentfulPlugin) TerraformRenderComponentDependsOn(site string, component string) []string {
	return []string{}
}

func (p *ContentfulPlugin) getSiteConfig(site string) *ContentfulConfig {
	cfg, ok := p.siteConfigs[site]
	if !ok {
		cfg = &ContentfulConfig{}
	}
	return cfg.extendConfig(p.globalConfig)
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
