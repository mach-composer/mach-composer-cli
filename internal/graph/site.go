package graph

import (
	"github.com/dominikbraun/graph"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type Site struct {
	baseNode
	NestedNodes []*SiteComponent
	SiteConfig  config.SiteConfig
}

func NewSite(g graph.Graph[string, Node], path, identifier string, deploymentType config.DeploymentType, ancestor Node,
	siteConfig config.SiteConfig) *Site {
	return &Site{
		baseNode:   newBaseNode(g, path, identifier, SiteType, ancestor, deploymentType),
		SiteConfig: siteConfig,
	}
}

func (s *Site) Hash() (string, error) {
	SortSiteComponentNodes(s.NestedNodes)

	var hashes []string
	for _, component := range s.NestedNodes {
		h, err := HashSiteComponentConfig(component.SiteComponentConfig)
		if err != nil {
			return "", err
		}
		hashes = append(hashes, h)
	}

	return utils.ComputeHash(hashes)
}
