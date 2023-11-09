package dependency

import "github.com/mach-composer/mach-composer-cli/internal/config"

type Project struct {
	node
	ProjectConfig *config.MachConfig
}
