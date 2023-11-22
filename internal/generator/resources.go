package generator

import (
	"embed"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

//go:embed templates/file_sources.tmpl
var fileSourcesTpl embed.FS

// renderFileSources uses templates/file_sources.tmpl to generate a terraform snippet for each file source
func renderFileSources(cfg *config.MachConfig, site *config.SiteConfig) (string, error) {
	tpl, err := fileSourcesTpl.ReadFile("templates/file_sources.tmpl")
	if err != nil {
		return "", err
	}
	return utils.RenderGoTemplate(string(tpl), cfg.Variables.GetEncryptedSources(site.Identifier))
}
