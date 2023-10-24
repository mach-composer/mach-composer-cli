package dependency

import "github.com/mach-composer/mach-composer-cli/internal/config"

type Project struct {
	nodeImpl

	Config *config.MachConfig
}
