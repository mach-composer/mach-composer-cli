package dependency

import (
	"github.com/mach-composer/mach-composer-cli/internal/config"
)

type Project struct {
	baseNode
	ProjectConfig *config.MachConfig
}

func (p *Project) Hash() (string, error) {
	return "", nil
}

func (p *Project) HasChanges() (bool, error) {
	return false, nil
}
