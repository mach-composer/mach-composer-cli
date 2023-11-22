package dependency

import (
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type Site struct {
	baseNode
	NestedSiteComponentConfigs []config.SiteComponentConfig
	SiteConfig                 config.SiteConfig
}

func (s *Site) Hash() (string, error) {
	var componentHashes []string
	for _, component := range s.NestedSiteComponentConfigs {
		hash, err := component.Hash()
		if err != nil {
			return "", err
		}
		componentHashes = append(componentHashes, hash)
	}

	return utils.ComputeHash(componentHashes)
}

func (s *Site) HasChanges() (bool, error) {
	hash, err := s.Hash()
	if err != nil {
		return true, err
	}

	var componentHashes []string
	for _, component := range s.NestedSiteComponentConfigs {
		tfOutput, err := utils.ParseSiteComponentOutputByKey(s.outputs, component.Name)
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

		componentHashes = append(componentHashes, *tfHash)
	}

	tfHash, err := utils.ComputeHash(componentHashes)
	if err != nil {
		return false, err
	}

	return hash != tfHash, nil
}
