package dependency

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/rs/zerolog/log"
)

type Site struct {
	node
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

func (s *Site) HasConfigChanges(ctx context.Context) (bool, error) {
	hash, err := s.Hash()
	if err != nil {
		return true, err
	}

	path := fmt.Sprintf("deployments/%s", s.Path())

	tfOutput, err := utils.GetTerraformOutput(ctx, path)
	if err != nil {
		return false, err
	}

	var componentHashes []string
	for _, component := range s.NestedSiteComponentConfigs {
		out, exists := tfOutput.GetSiteComponentOutput(component.Name)
		if !exists {
			log.Info().Msgf("Component %s does not exist in terraform output. Assuming update is required", component.Name)
			return true, nil
		}

		componentHashes = append(componentHashes, out.Value.Hash)
	}

	tfHash, err := utils.ComputeHash(componentHashes)
	if err != nil {
		return false, err
	}

	return hash != tfHash, nil
}
