package mcsdk

import (
	"github.com/hashicorp/go-hclog"
)

type PluginSchema struct {
	Identifier string

	Configure                   func(environment, provider string) error
	IsEnabled                   func() bool
	SetRemoteStateBackend       func(data map[string]any) error
	SetGlobalConfig             func(data map[string]any) error
	SetSiteConfig               func(site string, data map[string]any) error
	SetSiteComponentConfig      func(site string, component string, data map[string]any) error
	SetSiteEndpointsConfig      func(site string, data map[string]any) error
	SetComponentConfig          func(component string, data map[string]any) error
	SetComponentEndpointsConfig func(component string, endpoints map[string]string) error
	RenderTerraformStateBackend func(site string) (string, error)
	RenderTerraformProviders    func(site string) (string, error)
	RenderTerraformResources    func(site string) (string, error)
	RenderTerraformComponent    func(site string, component string) (*ComponentSnippets, error)
}

type Adapter struct {
	Logger hclog.Logger
	fn     *PluginSchema
}

func NewPlugin(fn *PluginSchema) MachComposerPlugin {
	return &Adapter{
		fn: fn,
	}
}

func (c *Adapter) SetLogger(logger hclog.Logger) {
	c.Logger = logger
}

func (c *Adapter) Configure(environment, provider string) error {
	if c.fn.Configure != nil {
		return c.fn.Configure(environment, provider)
	}
	return nil
}

func (c *Adapter) Identifier() string {
	return c.fn.Identifier
}

func (p *Adapter) IsEnabled() bool {
	if p.fn.IsEnabled != nil {
		return p.fn.IsEnabled()
	}
	return true
}

func (p *Adapter) SetRemoteStateBackend(data map[string]any) error {
	if p.fn.IsEnabled != nil && p.fn.SetRemoteStateBackend != nil {
		return p.fn.SetRemoteStateBackend(data)
	}
	return nil
}

func (p *Adapter) SetGlobalConfig(data map[string]any) error {
	if p.fn.IsEnabled != nil && p.fn.SetGlobalConfig != nil {
		return p.fn.SetGlobalConfig(data)
	}
	return nil
}
func (p *Adapter) SetSiteConfig(site string, data map[string]any) error {
	if p.fn.IsEnabled != nil && p.fn.SetSiteConfig != nil {
		return p.fn.SetSiteConfig(site, data)
	}
	return nil
}

func (p *Adapter) SetSiteComponentConfig(site string, component string, data map[string]any) error {
	if p.fn.IsEnabled != nil && p.fn.SetSiteComponentConfig != nil {
		return p.fn.SetSiteComponentConfig(site, component, data)
	}
	return nil
}

func (p *Adapter) SetSiteEndpointsConfig(site string, data map[string]any) error {
	if p.fn.IsEnabled != nil {
		return p.fn.SetSiteEndpointsConfig(site, data)
	}
	return nil
}

func (p *Adapter) SetComponentConfig(component string, data map[string]any) error {
	if p.fn.SetComponentConfig != nil {
		return p.fn.SetComponentConfig(component, data)
	}
	return nil
}

func (p *Adapter) SetComponentEndpointsConfig(component string, endpoints map[string]string) error {
	if p.fn.SetComponentEndpointsConfig != nil {
		return p.fn.SetComponentEndpointsConfig(component, endpoints)
	}
	return nil
}
func (p *Adapter) RenderTerraformStateBackend(site string) (string, error) {
	if p.fn.IsEnabled != nil {
		return p.fn.RenderTerraformStateBackend(site)
	}
	return "", nil
}

func (p *Adapter) RenderTerraformProviders(site string) (string, error) {
	if p.fn.IsEnabled != nil && p.fn.RenderTerraformProviders != nil {
		return p.fn.RenderTerraformProviders(site)
	}
	return "", nil
}

func (p *Adapter) RenderTerraformResources(site string) (string, error) {
	if p.fn.IsEnabled != nil && p.fn.RenderTerraformResources != nil {
		return p.fn.RenderTerraformResources(site)
	}
	return "", nil
}

func (p *Adapter) RenderTerraformComponent(site, component string) (*ComponentSnippets, error) {
	if p.fn.IsEnabled != nil && p.fn.RenderTerraformComponent != nil {
		return p.fn.RenderTerraformComponent(site, component)
	}
	return nil, nil
}
