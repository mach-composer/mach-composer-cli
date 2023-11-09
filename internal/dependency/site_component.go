package dependency

import "github.com/mach-composer/mach-composer-cli/internal/config"

type SiteComponent struct {
	node
	SiteConfig          config.SiteConfig
	SiteComponentConfig config.SiteComponentConfig
}
