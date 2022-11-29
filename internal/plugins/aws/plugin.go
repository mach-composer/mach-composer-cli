package aws

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/elliotchance/pie/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"

	"github.com/labd/mach-composer/internal/plugins/shared"
)

type AWSPlugin struct {
	environment      string
	provider         string
	remoteState      *AWSTFState
	siteConfigs      map[string]*SiteConfig
	componentConfigs map[string]ComponentConfig
	endpointsConfigs map[string]map[string]EndpointConfig
}

func NewAWSPlugin() *AWSPlugin {
	return &AWSPlugin{
		provider:         "3.74.1",
		siteConfigs:      map[string]*SiteConfig{},
		componentConfigs: map[string]ComponentConfig{},
		endpointsConfigs: map[string]map[string]EndpointConfig{},
	}
}

func (p *AWSPlugin) Configure(environment string, provider string) error {
	p.environment = environment
	if provider != "" {
		p.provider = provider
	}
	return nil
}

func (p *AWSPlugin) IsEnabled() bool {
	return len(p.siteConfigs) > 0
}

func (p *AWSPlugin) Identifier() string {
	return "aws"
}

func (p *AWSPlugin) SetRemoteStateBackend(data map[string]any) error {
	state := &AWSTFState{}
	if err := mapstructure.Decode(data, state); err != nil {
		return err
	}
	if err := defaults.Set(state); err != nil {
		return err
	}
	p.remoteState = state
	return nil
}

func (p *AWSPlugin) SetGlobalConfig(data map[string]any) error {
	return nil
}

func (p *AWSPlugin) SetSiteConfig(site string, data map[string]any) error {
	cfg := SiteConfig{}

	if val, ok := data["account_id"].(int); ok {
		data["account_id"] = fmt.Sprintf("%d", val)
	}

	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}

	if err := defaults.Set(&cfg); err != nil {
		return err
	}

	p.siteConfigs[site] = &cfg
	return nil
}

func (p *AWSPlugin) SetSiteComponentConfig(site string, component string, data map[string]any) error {
	return nil
}

func (p *AWSPlugin) SetSiteEndpointsConfig(site string, data map[string]any) error {
	configs := map[string]EndpointConfig{}
	for epId, epData := range data {
		cfg := EndpointConfig{}
		if url, ok := epData.(string); ok {
			cfg.URL = url
		} else {
			if mapData, ok := epData.(map[string]any); ok {
				if val, ok := mapData["aws"].(map[string]any); ok {
					logrus.Warnln("the aws node on the endpoint will be removed. Set the children directly in the endpoint")
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

func (p *AWSPlugin) SetComponentConfig(component string, data map[string]any) error {
	return nil
}

func (p *AWSPlugin) SetComponentEndpointsConfig(component string, endpoints map[string]string) error {
	cfg := ComponentConfig{
		Endpoints: endpoints,
	}
	p.componentConfigs[component] = cfg
	return nil
}

func (p *AWSPlugin) TerraformRenderStateBackend(site string) (string, error) {
	if p.remoteState == nil {
		return "", nil
	}

	templateContext := struct {
		State *AWSTFState
		Site  string
	}{
		State: p.remoteState,
		Site:  site,
	}

	template := `
	backend "s3" {
	  bucket         = "{{ .State.Bucket }}"
	  key            = "{{ .State.KeyPrefix}}/{{ .Site }}"
	  region         = "{{ .State.Region }}"
	  {{ if .State.RoleARN }}
	  role_arn       = "{{ .State.RoleARN }}"
	  {{ end }}
	  {{ if .State.LockTable }}
	  dynamodb_table = "{{ .State.LockTable }}"
	  {{ end }}
	  encrypt        = {{ .State.Encrypt }}
	}
	`
	return shared.RenderGoTemplate(template, templateContext)
}

func (p *AWSPlugin) TerraformRenderProviders(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	result := fmt.Sprintf(`
		aws = {
			version = "%s"
		}`, p.provider)
	return result, nil
}

func (p *AWSPlugin) TerraformRenderResources(site string) (string, error) {
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
				return "", fmt.Errorf("component requires undeclared endpoint: %s", external)
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

	content, err := renderResources(site, p.environment, cfg, pie.Values(activeEndpoints))
	if err != nil {
		return "", fmt.Errorf("failed to render resources: %w", err)
	}

	return content, nil
}

func (p *AWSPlugin) TerraformRenderComponentResources(site string, component string) (string, error) {
	return "", nil
}

func (p *AWSPlugin) TerraformRenderComponentVars(site string, component string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	componentCfg := p.componentConfigs[component]

	endpointNames := map[string]string{}
	for key, value := range componentCfg.Endpoints {
		endpointNames[shared.Slugify(key)] = shared.Slugify(value)
	}

	templateContext := struct {
		Site      *SiteConfig
		Endpoints map[string]string
	}{
		Site:      cfg,
		Endpoints: endpointNames,
	}

	template := `
		{{ range $cEndpoint, $sEndpoint := .Endpoints }}
		aws_endpoint_{{ $cEndpoint }} = {
			url = local.endpoint_url_{{ $sEndpoint }}
			api_gateway_id = aws_apigatewayv2_api.{{ $sEndpoint }}_gateway.id
			api_gateway_execution_arn = aws_apigatewayv2_api.{{ $sEndpoint }}_gateway.execution_arn
		}
		{{ end }}
	`
	return shared.RenderGoTemplate(template, templateContext)
}

func (p *AWSPlugin) TerraformRenderComponentProviders(site string, component string) ([]string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return []string{}, nil
	}
	providers := []string{"aws = aws"}
	for _, provider := range cfg.ExtraProviders {
		providers = append(providers, fmt.Sprintf("aws.%s = aws.%s", provider.Name, provider.Name))
	}
	return providers, nil
}

func (p *AWSPlugin) TerraformRenderComponentDependsOn(site string, component string) ([]string, error) {
	// This shouldn't be needed since we already pass the values to the component
	// make it automatically depend on that value
	result := []string{}

	componentCfg := p.componentConfigs[component]
	for _, value := range componentCfg.Endpoints {
		depends := fmt.Sprintf("aws_apigatewayv2_api.%s_gateway", shared.Slugify(value))
		result = append(result, depends)
	}
	return result, nil
}

func (p *AWSPlugin) getSiteConfig(site string) *SiteConfig {
	cfg, ok := p.siteConfigs[site]
	if !ok {
		return nil
	}
	return cfg
}
