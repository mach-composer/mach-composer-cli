package dependency

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type SiteComponent struct {
	baseNode
	SiteConfig          config.SiteConfig
	SiteComponentConfig config.SiteComponentConfig
}

func (sc *SiteComponent) Hash() (string, error) {
	return sc.SiteComponentConfig.Hash()
}

func (sc *SiteComponent) HasChanges(ctx context.Context) (bool, error) {
	hash, err := sc.Hash()
	if err != nil {
		return true, err
	}

	path := fmt.Sprintf("deployments/%s", sc.Path())

	tfOutput, err := utils.GetTerraformOutputByKey(ctx, path, sc.identifier)
	if err != nil {
		return false, err
	}

	if tfOutput == nil {
		return true, nil
	}

	return hash != tfOutput.Value.Hash, nil
}
