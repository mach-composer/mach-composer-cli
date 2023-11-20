package dependency

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/config"
)

type Project struct {
	node
	ProjectConfig *config.MachConfig
}

func (p *Project) Hash() (string, error) {
	return "", nil
}

func (p *Project) HasConfigChanges(context.Context) (bool, error) {
	return false, nil
}
