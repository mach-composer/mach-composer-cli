package dependency

import (
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

func (sc *SiteComponent) HasChanges() (bool, error) {
	hash, err := sc.Hash()
	if err != nil {
		return true, err
	}

	tfOutput, err := utils.ParseSiteComponentOutputByKey(sc.outputs, sc.identifier)
	if err != nil {
		return false, err
	}

	if tfOutput == nil {
		return true, nil
	}

	tfHash := tfOutput.Value.Hash
	if tfHash == nil {
		return true, nil
	}

	return hash != *tfHash, nil
}
