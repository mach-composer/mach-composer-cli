package generator

import (
	"embed"
	"fmt"
	"strings"

	"github.com/elliotchance/pie/v2"
	"github.com/flosch/pongo2/v5"

	"github.com/labd/mach-composer/internal/config"
	"github.com/labd/mach-composer/internal/plugins/shared"
	"github.com/labd/mach-composer/internal/variables"
)

//go:embed templates/*
var templates embed.FS

type TemplateRenderer struct {
	templateSet *pongo2.TemplateSet

	componentTemplate *pongo2.Template
}

var renderer TemplateRenderer

func init() {
	renderer.templateSet = pongo2.NewSet("", &shared.EmbedLoader{Content: templates})
	renderer.componentTemplate = pongo2.Must(renderer.templateSet.FromFile("component.tf"))
}

// renderSite is responsible for generating the `site.tf` file. Therefore it is
// the main entrypoint for generating the terraform file for each site.
func renderSite(cfg *config.MachConfig, site *config.Site) (string, error) {
	result := []string{
		"# This file is auto-generated by MACH composer",
		fmt.Sprintf("# Site: %s", site.Identifier),
	}

	// Render the terraform config
	val, err := renderTerraformConfig(cfg, site)
	if err != nil {
		return "", err
	}
	result = append(result, val)

	// Render all the global resources
	val, err = renderTerraformResources(cfg, site)
	if err != nil {
		return "", err
	}
	result = append(result, val)

	// Render every component (incuding component specific resources)
	for i := range site.Components {
		component := &site.Components[i]
		val, err := renderComponent(cfg, site, component)
		if err != nil {
			return "", fmt.Errorf("failed to render component %s: %w", component.Name, err)
		}
		result = append(result, val)
	}

	content := strings.Join(result, "\n")
	return content, nil
}

// renderTerraformConfig renders the terraform settings block which defines the
// remote state to be used and the providers to be loaded.
func renderTerraformConfig(cfg *config.MachConfig, site *config.Site) (string, error) {
	providers := []string{}
	for _, plugin := range cfg.Plugins.All() {
		content, err := plugin.TerraformRenderProviders(site.Identifier)
		if err != nil {
			return "", err
		}
		providers = append(providers, content)
	}

	statePlugin, err := cfg.Plugins.Get(cfg.Global.TerraformStateProvider)
	if err != nil {
		return "", fmt.Errorf("failed to resolve plugin for terraform state: %w", err)
	}

	backendConfig, err := statePlugin.TerraformRenderStateBackend(site.Identifier)
	if err != nil {
		return "", fmt.Errorf("failed to render backend config: %w", err)
	}

	template := `
		terraform {
			{{ .BackendConfig }}

			required_providers {
				{{ range $provider := .Providers }}
					{{ $provider }}
				{{ end }}

				{{ if .IncludeSOPS }}
				sops = {
					source = "carlpett/sops"
					version = "~> 0.5"
				}
				{{ end }}
			}
	  }
	`
	templateContext := struct {
		Providers     []string
		BackendConfig string
		IncludeSOPS   bool
	}{
		Providers:     providers,
		BackendConfig: backendConfig,
		IncludeSOPS:   cfg.Variables.Encrypted,
	}
	return shared.RenderGoTemplate(template, templateContext)
}

func renderTerraformResources(cfg *config.MachConfig, site *config.Site) (string, error) {
	resources := []string{}
	for _, plugin := range cfg.Plugins.All() {
		content, err := plugin.TerraformRenderResources(site.Identifier)
		if err != nil {
			return "", err
		}
		resources = append(resources, content)
	}

	template := `
		{{ if .VarsEncrypted }}
			data "local_file" "variables" {
			filename = "{{ .VarsFilename }}"
		}

		data "sops_external" "variables" {
			source     = data.local_file.variables.content
			input_type = "yaml"
		}
		{{ end }}

		# Plugins
		{{ range $resource := .Resources }}
			{{ $resource }}
		{{ end }}
	`
	templateContext := struct {
		Resources     []string
		VarsFilename  string
		VarsEncrypted bool
	}{
		Resources:     resources,
		VarsFilename:  cfg.Variables.Filepath,
		VarsEncrypted: cfg.Variables.Encrypted,
	}
	return shared.RenderGoTemplate(template, templateContext)
}

// renderComponent uses templates/component.tf to generate a terraform snippet
// for each component
func renderComponent(cfg *config.MachConfig, site *config.Site, component *config.SiteComponent) (string, error) {
	pVars := []string{}
	pResources := []string{}
	pDependsOn := []string{}
	pProviders := []string{}
	for _, plugin := range cfg.Plugins.All() {
		if !pie.Contains(component.Definition.Integrations, plugin.Identifier()) {
			continue
		}

		value, err := plugin.TerraformRenderComponentResources(site.Identifier, component.Name)
		if err != nil {
			return "", err
		}
		pResources = append(pResources, value)

		value, err = plugin.TerraformRenderComponentVars(site.Identifier, component.Name)
		if err != nil {
			return "", err
		}
		pVars = append(pVars, value)

		values, err := plugin.TerraformRenderComponentProviders(site.Identifier, component.Name)
		if err != nil {
			return "", err
		}
		pProviders = append(pProviders, values...)

		values, err = plugin.TerraformRenderComponentDependsOn(site.Identifier, component.Name)
		if err != nil {
			return "", err
		}
		pDependsOn = append(pDependsOn, values...)
	}

	componentVariables, err := variables.InterpolateComponentVars(component.Variables)
	if err != nil {
		return "", err
	}

	componentSecrets, err := variables.InterpolateComponentVars(component.Secrets)
	if err != nil {
		return "", err
	}

	return renderer.componentTemplate.Execute(pongo2.Context{
		"global":             cfg.Global,
		"site":               site,
		"variables":          cfg.Variables,
		"component":          component,
		"componentSecrets":   componentSecrets,
		"componentVariables": componentVariables,
		"definition":         component.Definition,
		"pluginVariables":    pVars,
		"pluginResources":    pResources,
		"pluginProviders":    pProviders,
		"pluginDependsOn":    pDependsOn,
	})
}
