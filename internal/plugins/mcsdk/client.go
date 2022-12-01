package mcsdk

import (
	"net/rpc"

	"github.com/sirupsen/logrus"
)

type PluginRPC struct {
	client     *rpc.Client
	identifier string
}

func (p *PluginRPC) Identifier() string {
	// return p.identifier
	var resp string
	err := p.client.Call("Plugin.Identifier", new(any), &resp)
	if err != nil {
		logrus.Error(err)
		return ""
	}
	return resp
}

func (p *PluginRPC) IsEnabled() bool {
	var resp bool
	err := p.client.Call("Plugin.IsEnabled", new(any), &resp)
	if err != nil {
		logrus.Error(err)
		return false
	}
	return resp
}

func (p *PluginRPC) Configure(environment string, provider string) error {
	param := ConfigureInput{
		Environment: environment,
		Provider:    provider,
	}
	resp := ErrorOutput{}
	err := p.client.Call("Plugin.Configure", param, &resp)
	if err != nil {
		return err
	}
	return resp.Err
}

func (p *PluginRPC) SetRemoteStateBackend(data map[string]any) error {
	param := SetRemoteStateBackendInput{
		Data: data,
	}
	resp := ErrorOutput{}
	err := p.client.Call("Plugin.SetRemoteStateBackend", param, &resp)
	if err != nil {
		return err
	}
	return resp.Err
}

func (p *PluginRPC) SetGlobalConfig(data map[string]any) error {
	param := SetGlobalConfigInput{
		Data: data,
	}
	resp := ErrorOutput{}
	err := p.client.Call("Plugin.SetGlobalConfig", param, &resp)
	if err != nil {
		return err
	}
	return resp.Err
}

func (p *PluginRPC) SetSiteConfig(name string, data map[string]any) error {
	param := SetSiteConfigInput{
		Name: name,
		Data: data,
	}
	resp := ErrorOutput{}
	err := p.client.Call("Plugin.SetSiteConfig", param, &resp)
	if err != nil {
		return err
	}
	return resp.Err
}

func (p *PluginRPC) SetSiteComponentConfig(site string, component string, data map[string]any) error {
	param := SetSiteComponentConfigInput{
		Site:      site,
		Component: component,
		Data:      data,
	}
	resp := SetSiteComponentConfigOutput{}
	err := p.client.Call("Plugin.SetSiteComponentConfig", param, &resp)
	if err != nil {
		return err
	}
	return resp.Err
}

func (p *PluginRPC) SetSiteEndpointsConfig(site string, data map[string]any) error {
	param := SetSiteEndpointsConfigInput{
		Site: site,
		Data: data,
	}
	resp := SetSiteEndpointsConfigOutput{}
	err := p.client.Call("Plugin.SetSiteEndpointsConfig", param, &resp)
	if err != nil {
		return err
	}
	return resp.Err
}

func (p *PluginRPC) SetComponentConfig(component string, data map[string]any) error {
	param := SetComponentConfigInput{
		Component: component,
		Data:      data,
	}
	resp := SetComponentConfigOutput{}
	err := p.client.Call("Plugin.SetComponentConfig", param, &resp)
	if err != nil {
		return err
	}
	return resp.Err
}

func (p *PluginRPC) SetComponentEndpointsConfig(component string, endpoints map[string]string) error {
	param := SetComponentEndpointsConfigInput{
		Component: component,
		Endpoints: endpoints,
	}
	resp := SetComponentEndpointsConfigOutput{}
	err := p.client.Call("Plugin.SetComponentEndpointsConfig", param, &resp)
	if err != nil {
		return err
	}
	return resp.Err

}

func (p *PluginRPC) RenderTerraformStateBackend(site string) (string, error) {
	param := RenderTerraformStateBackendInput{
		Site: site,
	}
	resp := RenderTerraformStateBackendOutput{}
	err := p.client.Call("Plugin.RenderTerraformStateBackend", param, &resp)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return resp.Result, resp.Err
}

func (p *PluginRPC) RenderTerraformProviders(site string) (string, error) {
	param := RenderTerraformProvidersInput{
		Site: site,
	}
	resp := RenderTerraformProvidersOutput{}
	err := p.client.Call("Plugin.RenderTerraformProviders", param, &resp)
	if err != nil {
		return "", err
	}
	return resp.Result, resp.Err
}

func (p *PluginRPC) RenderTerraformResources(site string) (string, error) {
	param := RenderTerraformResourcesInput{
		Site: site,
	}
	resp := RenderTerraformResourcesOutput{}
	err := p.client.Call("Plugin.RenderTerraformResources", param, &resp)
	if err != nil {
		return "", err
	}
	return resp.Result, resp.Err
}

func (p *PluginRPC) RenderTerraformComponent(site string, component string) (*ComponentSnippets, error) {
	param := RenderTerraformComponentInput{
		Site:      site,
		Component: component,
	}
	resp := RenderTerraformComponentOutput{}
	err := p.client.Call("Plugin.RenderTerraformComponent", param, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Result, resp.Err
}
