package plugins

import (
	schemav1 "github.com/mach-composer/mach-composer-plugin-sdk/schema"
	schemav2 "github.com/mach-composer/mach-composer-plugin-sdk/v2/schema"
)

// PluginV1Adapter will adapt the old v1 plugin to work with the new v2 settings
type PluginV1Adapter struct {
	v1Plugin schemav1.MachComposerPlugin
}

func NewPluginV1Adapter(p schemav1.MachComposerPlugin) *PluginV1Adapter {
	return &PluginV1Adapter{v1Plugin: p}
}

func (p *PluginV1Adapter) Identifier() string {
	return p.v1Plugin.Identifier()
}

func (p *PluginV1Adapter) Configure(environment string, provider string) error {
	return p.v1Plugin.Configure(environment, provider)
}

func (p *PluginV1Adapter) GetValidationSchema() (*schemav2.ValidationSchema, error) {
	v1ValidationSchema, err := p.v1Plugin.GetValidationSchema()
	if err != nil {
		return nil, err
	}
	return &schemav2.ValidationSchema{
		GlobalConfigSchema:             v1ValidationSchema.GlobalConfigSchema,
		RemoteStateSchema:              v1ValidationSchema.RemoteStateSchema,
		SiteConfigSchema:               v1ValidationSchema.SiteConfigSchema,
		SiteComponentConfigSchema:      v1ValidationSchema.SiteComponentConfigSchema,
		SiteEndpointConfig:             v1ValidationSchema.SiteEndpointConfig,
		ComponentConfigSchema:          v1ValidationSchema.ComponentConfigSchema,
		ComponentEndpointsConfigSchema: v1ValidationSchema.ComponentEndpointsConfigSchema,
	}, nil
}

func (p *PluginV1Adapter) SetGlobalConfig(data map[string]any) error {
	return p.v1Plugin.SetGlobalConfig(data)
}

func (p *PluginV1Adapter) SetSiteConfig(site string, data map[string]any) error {
	return p.v1Plugin.SetSiteConfig(site, data)
}

func (p *PluginV1Adapter) SetSiteComponentConfig(site string, component string, data map[string]any) error {
	return p.v1Plugin.SetSiteComponentConfig(site, component, data)
}

func (p *PluginV1Adapter) SetSiteEndpointConfig(site string, name string, data map[string]any) error {
	return p.v1Plugin.SetSiteEndpointConfig(site, name, data)
}

func (p *PluginV1Adapter) SetComponentConfig(component string, _ string, data map[string]any) error {
	return p.v1Plugin.SetComponentConfig(component, data)
}

func (p *PluginV1Adapter) SetComponentEndpointsConfig(component string, data map[string]string) error {
	return p.v1Plugin.SetComponentEndpointsConfig(component, data)
}

func (p *PluginV1Adapter) RenderTerraformProviders(site string) (string, error) {
	return p.v1Plugin.RenderTerraformProviders(site)
}

func (p *PluginV1Adapter) RenderTerraformResources(site string) (string, error) {
	return p.v1Plugin.RenderTerraformResources(site)
}

func (p *PluginV1Adapter) RenderTerraformComponent(site string, component string) (*schemav2.ComponentSchema, error) {
	v1ComponentSchema, err := p.v1Plugin.RenderTerraformComponent(site, component)
	if err != nil {
		return nil, err
	}

	return &schemav2.ComponentSchema{
		Resources: v1ComponentSchema.Resources,
		Variables: v1ComponentSchema.Variables,
		DependsOn: v1ComponentSchema.DependsOn,
		Providers: v1ComponentSchema.Providers,
	}, nil
}
