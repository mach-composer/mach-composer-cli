package generator

import (
	"embed"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

//go:embed templates/terraform.tmpl
var terraformTpl embed.FS

func renderTerraformConfig(cfg *config.MachConfig, site *config.SiteConfig, target string) (string, error) {
	var providers []string
	for _, plugin := range cfg.Plugins.All() {
		content, err := plugin.RenderTerraformProviders(site.Identifier)
		if err != nil {
			return "", fmt.Errorf("plugin %s failed to render providers: %w", plugin.Name, err)
		}
		if content != "" {
			providers = append(providers, content)
		}
	}

	if !cfg.StateRepository.Has(target) {
		return "", fmt.Errorf("state repository does not have a backend for site %s", target)
	}
	backendConfig, err := cfg.StateRepository.Get(target).Backend()
	if err != nil {
		return "", err
	}

	tpl, err := terraformTpl.ReadFile("templates/terraform.tmpl")
	if err != nil {
		return "", err
	}

	templateContext := struct {
		Providers     []string
		BackendConfig string
		RemoteState   string
		IncludeSOPS   bool
	}{
		Providers:     providers,
		BackendConfig: backendConfig,
		IncludeSOPS:   cfg.Variables.HasEncrypted(site.Identifier),
	}
	return utils.RenderGoTemplate(string(tpl), templateContext)
}
