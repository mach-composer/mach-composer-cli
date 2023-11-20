package dependency

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type SiteComponent struct {
	node
	SiteConfig          config.SiteConfig
	SiteComponentConfig config.SiteComponentConfig
}

func (sc *SiteComponent) Hash() (string, error) {
	return sc.SiteComponentConfig.Hash()
}

func (sc *SiteComponent) HasConfigChanges(ctx context.Context) (bool, error) {
	hash, err := sc.Hash()
	if err != nil {
		return true, err
	}

	path := fmt.Sprintf("deployments/%s", sc.Path())

	tfOutput, err := utils.GetTerraformOutput(ctx, path)
	if err != nil {
		return false, err
	}

	tfHash, exists := tfOutput.GetSiteComponentOutput(sc.SiteComponentConfig.Name)
	if !exists {
		return true, nil
	}

	return hash != tfHash.Value.Hash, nil
}
