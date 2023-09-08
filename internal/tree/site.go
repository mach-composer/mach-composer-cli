package tree

import (
	"github.com/mach-composer/mach-composer-cli/internal/config"
)

type Site struct {
	nodeImpl
	Config *config.SiteConfig
}
