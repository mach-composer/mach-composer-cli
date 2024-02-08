package graph

import (
	"github.com/dominikbraun/graph"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/rs/zerolog/log"
	"sort"
)

type Site struct {
	baseNode
	NestedSiteComponentConfigs []config.SiteComponentConfig
	SiteConfig                 config.SiteConfig
}

func NewSite(g graph.Graph[string, Node], path, identifier string, deploymentType config.DeploymentType, ancestor Node,
	siteConfig config.SiteConfig) *Site {
	return &Site{
		baseNode:   newBaseNode(g, path, identifier, SiteType, ancestor, deploymentType),
		SiteConfig: siteConfig,
	}
}

func (s *Site) Hash() (string, error) {
	sort.Slice(s.NestedSiteComponentConfigs, func(i, j int) bool {
		return s.NestedSiteComponentConfigs[i].Name < s.NestedSiteComponentConfigs[j].Name
	})

	var componentHashes []string
	for _, component := range s.NestedSiteComponentConfigs {
		hash, err := hashSiteComponentConfig(component)
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
		return false, err
	}

	tfHash, err := utils.ParseHashOutput(s.outputs)
	if err != nil {
		log.Warn().Msgf("Could not parse hash output: %s", err)
		return true, nil
	}

	return hash != tfHash, nil
}