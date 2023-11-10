package dependency

import "github.com/mach-composer/mach-composer-cli/internal/config"

type Site struct {
	node
	NestedSiteComponentConfigs []config.SiteComponentConfig
	SiteConfig                 config.SiteConfig
}
