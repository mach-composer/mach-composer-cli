package generator

import (
	"embed"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

//go:embed templates/resources.tmpl
var resourcesTpl embed.FS

func renderResources(cfg *config.MachConfig, site *config.SiteConfig) (string, error) {
	var resources []string
	for _, plugin := range cfg.Plugins.All() {
		content, err := plugin.RenderTerraformResources(site.Identifier)
		if err != nil {
			return "", fmt.Errorf("plugin %s failed to render resources: %w", plugin.Name, err)
		}

		if content != "" {
			resources = append(resources, content)
		}
	}

	tpl, err := resourcesTpl.ReadFile("templates/resources.tmpl")
	if err != nil {
		return "", err
	}

	return utils.RenderGoTemplate(string(tpl), resources)
}

//go:embed templates/file_sources.tmpl
var fileSourcesTpl embed.FS

func renderFileSources(cfg *config.MachConfig, site *config.SiteConfig) (string, error) {
	tpl, err := fileSourcesTpl.ReadFile("templates/file_sources.tmpl")
	if err != nil {
		return "", err
	}
	return utils.RenderGoTemplate(string(tpl), cfg.Variables.GetEncryptedSources(site.Identifier))
}
