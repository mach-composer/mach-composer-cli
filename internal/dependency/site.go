package dependency

import "github.com/mach-composer/mach-composer-cli/internal/config"

type Site struct {
	nodeImpl

	ProjectConfig *config.MachConfig
	Config        *config.SiteConfig
}
