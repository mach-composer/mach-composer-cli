package azure

import (
	"fmt"
	"log"
	"os"

	"github.com/creasty/defaults"
	"github.com/elliotchance/pie/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"

	"github.com/labd/mach-composer/internal/plugins/shared"
)

type AzurePlugin struct {
	environment      string
	remoteState      *AzureTFState
	globalConfig     *GlobalConfig
	siteConfigs      map[string]SiteConfig
	componentConfigs map[string]ComponentConfig
	endpointsConfigs map[string]map[string]EndpointConfig
}

func NewAzurePlugin() *AzurePlugin {
	return &AzurePlugin{
		siteConfigs:      map[string]SiteConfig{},
		componentConfigs: map[string]ComponentConfig{},
		endpointsConfigs: map[string]map[string]EndpointConfig{},
	}
}

func (p *AzurePlugin) Initialize(environment string) error {
	p.environment = environment
	return nil
}

func (p *AzurePlugin) IsEnabled() bool {
	return len(p.siteConfigs) > 0
}

func (p *AzurePlugin) Identifier() string {
	return "azure"
}

func (p *AzurePlugin) SetRemoteStateBackend(data map[string]any) error {
	state := &AzureTFState{}
	if err := mapstructure.Decode(data, state); err != nil {
		return err
	}
	if err := defaults.Set(state); err != nil {
		return err
	}
	p.remoteState = state
	return nil
}

func (p *AzurePlugin) SetGlobalConfig(data map[string]any) error {
	if err := mapstructure.Decode(data, &p.globalConfig); err != nil {
		return err
	}
	return nil
}

func (p *AzurePlugin) SetSiteConfig(site string, data map[string]any) error {
	cfg := SiteConfig{}

	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	cfg.merge(p.globalConfig)

	if cfg.ResourceGroup != "" {
		fmt.Fprintf(
			os.Stderr,
			"WARNING: resource_group on %s is used (%s). "+
				"Make sure it wasn't managed by MACH before otherwise "+
				"the resource group will get deleted.",
			site, cfg.ResourceGroup,
		)
	}

	p.siteConfigs[site] = cfg
	return nil
}

func (p *AzurePlugin) SetSiteComponentConfig(site string, component string, data map[string]any) error {
	return nil
}

func (p *AzurePlugin) SetSiteEndpointsConfig(site string, data map[string]any) error {
	configs := map[string]EndpointConfig{}
	for epId, epData := range data {
		cfg := EndpointConfig{}
		if url, ok := epData.(string); ok {
			cfg.URL = url
		} else {
			if mapData, ok := epData.(map[string]any); ok {
				if val, ok := mapData["azure"].(map[string]any); ok {
					logrus.Warnln("Warning: the azure node on the endpoint will be removed. Set the children directly in the endpoint")
					for key, value := range val {
						mapData[key] = value
					}
				}
			}

			if err := mapstructure.Decode(epData, &cfg); err != nil {
				return err
			}
		}

		if err := defaults.Set(&cfg); err != nil {
			return err
		}

		cfg.Key = epId
		configs[epId] = cfg
	}

	p.endpointsConfigs[site] = configs
	return nil
}

func (p *AzurePlugin) SetComponentConfig(component string, data map[string]any) error {
	cfg, ok := p.componentConfigs[component]
	if !ok {
		cfg = ComponentConfig{}
	}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	p.componentConfigs[component] = cfg
	return nil
}

func (p *AzurePlugin) SetComponentEndpointsConfig(component string, endpoints map[string]string) error {
	cfg, ok := p.componentConfigs[component]
	if ok {
		cfg.Endpoints = endpoints
	}
	return nil
}

func (p *AzurePlugin) TerraformRenderStateBackend(site string) (string, error) {
	templateContext := struct {
		State *AzureTFState
		Site  string
	}{
		State: p.remoteState,
		Site:  site,
	}

	template := `
	backend "azurerm" {
	  resource_group_name  = "{{ .State.ResourceGroup }}"
	  storage_account_name = "{{ .State.StorageAccount }}"
	  container_name       = "{{ .State.ContainerName }}"
	  key                  = "{{ .State.StateFolder}}/{{ .Site }}"
	}
	`
	return shared.RenderGoTemplate(template, templateContext)
}

func (p *AzurePlugin) TerraformRenderProviders(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	return `
    azure = {
      version = "2.99.0"
    }
	`, nil
}

func (p *AzurePlugin) TerraformRenderResources(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	activeEndpoints := map[string]EndpointConfig{}
	siteEndpoint := p.endpointsConfigs[site]

	needsDefaultEndpoint := false
	for _, component := range p.componentConfigs {
		for _, external := range component.Endpoints {
			if external == "default" {
				needsDefaultEndpoint = true
			}

			endpointConfig, ok := siteEndpoint[external]
			if !ok && external != "default" {
				log.Fatalf("component requires undeclared endpoint: %s", external)
			}

			if _, ok := activeEndpoints[external]; !ok {
				activeEndpoints[external] = endpointConfig
			}
		}
	}

	if needsDefaultEndpoint {
		activeEndpoints["default"] = EndpointConfig{
			Key: "default",
		}
	}

	return renderResources(site, p.environment, cfg, pie.Values(activeEndpoints))
}

func (p *AzurePlugin) TerraformRenderComponentResources(site string, component string) (string, error) {
	return "", nil
}

func (p *AzurePlugin) TerraformRenderComponentVars(site string, component string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	componentConfig, ok := p.componentConfigs[component]
	if !ok {
		componentConfig = ComponentConfig{} // TODO
	}

	endpointNames := map[string]string{}
	for key, value := range componentConfig.Endpoints {
		endpointNames[shared.Slugify(key)] = shared.Slugify(value)
	}

	templateContext := struct {
		Config      *SiteConfig
		Component   *ComponentConfig
		ServicePlan string
		Endpoints   map[string]string
	}{
		Config:      cfg,
		Component:   &componentConfig,
		ServicePlan: azureServicePlanResourceName(componentConfig.ServicePlan),
		Endpoints:   endpointNames,
	}

	template := `
		### azure related
		azure_short_name              = "{{ .Component.ShortName }}"
		azure_name_prefix             = local.name_prefix
		azure_subscription_id         = local.subscription_id
		azure_tenant_id               = local.tenant_id
		azure_region                  = local.region
		azure_service_object_ids      = local.service_object_ids
		azure_resource_group          = {
			name     = local.resource_group_name
			location = local.resource_group_location
		}

		{{ if .ServicePlan }}
		azure_app_service_plan = {
			id                  = azurerm_app_service_plan.{{ .ServicePlan }}.id
			name                = azurerm_app_service_plan.{{ .ServicePlan }}.name
			resource_group_name = azurerm_app_service_plan.{{ .ServicePlan }}.resource_group_name
		}
		{{ end }}

		{{ if .Config.AlertGroup }}
		azure_monitor_action_group_id = azurerm_monitor_action_group.alert_action_group.id
		{{ end }}

		{{ range $cEndpoint, $sEndpoint := .Endpoints }}
		azure_endpoint_{{ $cEndpoint }} = {
			url = local.endpoint_url_{{ $sEndpoint }}
			frontdoor_id = azurerm_frontdoor.app-service.header_frontdoor_id
		}
		{{ end }}
	`
	return shared.RenderGoTemplate(template, templateContext)
}

func (p *AzurePlugin) TerraformRenderComponentProviders(site string, component string) ([]string, error) {
	return []string{"azurerm = azurerm"}, nil
}

func (p *AzurePlugin) TerraformRenderComponentDependsOn(site string, component string) ([]string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return []string{}, nil
	}
	componentCfg := p.getComponentConfig(component)

	if componentCfg.ServicePlan != "" {
		if componentCfg.ServicePlan == "default" {
			return []string{"azurerm_app_service_plan.functionapps"}, nil
		} else {
			val := fmt.Sprintf("azurerm_app_service_plan.functionapps_%s", componentCfg.ServicePlan)
			return []string{val}, nil
		}
	}
	return []string{}, nil
}

func (p *AzurePlugin) getSiteConfig(site string) *SiteConfig {
	cfg, ok := p.siteConfigs[site]
	if !ok {
		return nil
	}
	return &cfg
}

func (p *AzurePlugin) getComponentConfig(name string) *ComponentConfig {
	componentConfig, ok := p.componentConfigs[name]
	if !ok {
		componentConfig = ComponentConfig{} // TODO
	}
	return &componentConfig
}
