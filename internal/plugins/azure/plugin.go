package azure

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/creasty/defaults"
	"github.com/elliotchance/pie/v2"
	"github.com/mitchellh/mapstructure"
)

type AzurePlugin struct {
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
					fmt.Println("Warning: the azure node on the endpoint will be removed. Set the children directly in the endpoint")
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

func (p *AzurePlugin) getSiteConfig(site string) *SiteConfig {
	cfg, ok := p.siteConfigs[site]
	if !ok {
		return nil
	}
	return &cfg
}

func (p *AzurePlugin) TerraformRenderStateBackend(site string) string {
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
	return renderTemplate(template, templateContext)
}

func (p *AzurePlugin) TerraformRenderProviders(site string) string {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return ""
	}

	return `
    azure = {
      version = "2.99.0"
    }
	`
}

func (p *AzurePlugin) TerraformRenderResources(site string) string {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return ""
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

	content, err := renderResources(site, cfg, pie.Values(activeEndpoints))
	if err != nil {
		panic(err)
	}

	return content
}

func (p *AzurePlugin) TerraformRenderComponentResources(site string, component string) string {
	return ""
}

func (p *AzurePlugin) TerraformRenderComponentVars(site string, component string) string {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return ""
	}

	componentConfig, ok := p.componentConfigs[component]
	if !ok {
		componentConfig = ComponentConfig{} // TODO
	}

	templateContext := struct {
		Config      *SiteConfig
		Component   *ComponentConfig
		ServicePlan string
	}{
		Config:      cfg,
		Component:   &componentConfig,
		ServicePlan: azureServicePlanResourceName(componentConfig.ServicePlan),
	}

	// {% for component_endpoint, site_endpoint in component.Endpoints -%}
	// azure_endpoint_{{ component_endpoint|slugify }} = {
	//   url = local.endpoint_url_{{ site_endpoint|slugify }}
	//   frontdoor_id = azurerm_frontdoor.app-service.header_frontdoor_id
	// }
	// {% endfor %}

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
	azure_app_service_plan        = {
	  id                  = azurerm_app_service_plan.{{ .ServicePlan }}.id
	  name                = azurerm_app_service_plan.{{ .ServicePlan }}.name
	  resource_group_name = azurerm_app_service_plan.{{ .ServicePlan }}.resource_group_name
	}
	{{ end }}
	{{ if .Config.AlertGroup }}
	azure_monitor_action_group_id = azurerm_monitor_action_group.alert_action_group.id
	{{ end }}
	`
	return renderTemplate(template, templateContext)
}

func (p *AzurePlugin) TerraformRenderComponentProviders(site string, component string) []string {
	return []string{"azurerm = azurerm"}
}

func (p *AzurePlugin) TerraformRenderComponentDependsOn(site string, component string) []string {
	return []string{"null_resource.commercetools"}
	// {% if site.Azure and component.Azure.ServicePlan %}
	// {% if component.Azure.ServicePlan == "default" %}
	// azurerm_app_service_plan.functionapps,{% else %}
	// azurerm_app_service_plan.functionapps_{{ component.Azure.ServicePlan }},{% endif %}
	// {% endif %}
}

func renderTemplate(t string, data any) string {
	tpl, err := template.New("template-1").Parse(t)
	if err != nil {
		panic(err)
	}

	var content bytes.Buffer
	if err := tpl.Execute(&content, data); err != nil {
		panic(err)
	}
	return content.String()
}
