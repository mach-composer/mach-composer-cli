package contentful

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/labd/mach-composer/internal/plugins/mcsdk"
	"github.com/labd/mach-composer/internal/plugins/shared"
)

type ContentfulPlugin struct {
	environment  string
	provider     string
	globalConfig *ContentfulConfig
	siteConfigs  map[string]*ContentfulConfig
	enabled      bool
}

func NewContentfulPlugin() mcsdk.MachComposerPlugin {
	state := &ContentfulPlugin{
		provider:    "0.1.0",
		siteConfigs: map[string]*ContentfulConfig{},
	}
	return mcsdk.NewPlugin(&mcsdk.PluginSchema{
		Identifier: "contentful",

		Configure: state.Configure,
		IsEnabled: state.IsEnabled,

		// Config
		SetGlobalConfig: state.SetGlobalConfig,
		SetSiteConfig:   state.SetSiteConfig,

		// Renders
		RenderTerraformProviders: state.TerraformRenderProviders,
		RenderTerraformResources: state.TerraformRenderResources,
		RenderTerraformComponent: state.RenderTerraformComponent,
	})
}

func (p *ContentfulPlugin) Configure(environment string, provider string) error {
	p.environment = environment
	if provider != "" {
		p.provider = provider
	}
	return nil
}

func (p *ContentfulPlugin) IsEnabled() bool {
	return p.enabled
}

func (p *ContentfulPlugin) Identifier() string {
	return "contentful"
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

func (p *ContentfulPlugin) TerraformRenderStateBackend(site string) (string, error) {
	return "", nil
}

func (p *ContentfulPlugin) TerraformRenderProviders(site string) (string, error) {
	result := fmt.Sprintf(`
		contentful = {
			source = "labd/contentful"
			version = "%s"
		}`, shared.VersionConstraint(p.provider))
	return result, nil
}

func (p *ContentfulPlugin) TerraformRenderResources(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
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
	return shared.RenderGoTemplate(template, cfg)
}

func (p *ContentfulPlugin) RenderTerraformComponent(site string, component string) (*mcsdk.ComponentSnippets, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return nil, nil
	}

	result := &mcsdk.ComponentSnippets{
		Variables: "contentful_space_id = contentful_space.space.id",
	}
	return result, nil
}

func (p *ContentfulPlugin) getSiteConfig(site string) *ContentfulConfig {
	cfg, ok := p.siteConfigs[site]
	if !ok {
		cfg = &ContentfulConfig{}
	}
	return cfg.extendConfig(p.globalConfig)
}
