package dependency

import "github.com/mach-composer/mach-composer-cli/internal/config"

type SiteComponent struct {
	nodeImpl

	ProjectConfig *config.MachConfig
	SiteConfig    *config.SiteConfig
	Config        *config.SiteComponentConfig
}
