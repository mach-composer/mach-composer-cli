package tree

import (
	"github.com/mach-composer/mach-composer-cli/internal/config"
)

type SiteComponent struct {
	nodeImpl
	Config *config.SiteComponentConfig
}
