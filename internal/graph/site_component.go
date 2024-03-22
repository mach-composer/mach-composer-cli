package graph

import (
	"github.com/dominikbraun/graph"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"sort"
)

type SiteComponent struct {
	baseNode
	SiteConfig          config.SiteConfig
	SiteComponentConfig config.SiteComponentConfig
}

func NewSiteComponent(g graph.Graph[string, Node], path, identifier string, deploymentType config.DeploymentType,
	ancestor Node, siteConfig config.SiteConfig, siteComponentConfig config.SiteComponentConfig) *SiteComponent {
	return &SiteComponent{
		baseNode:   newBaseNode(g, path, identifier, SiteComponentType, ancestor, deploymentType),
		SiteConfig: siteConfig, SiteComponentConfig: siteComponentConfig,
	}
}

func (sc *SiteComponent) Hash() (string, error) {
	return HashSiteComponentConfig(sc.SiteComponentConfig)
}

func SortSiteComponentNodes(nodes []*SiteComponent) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].SiteComponentConfig.Name < nodes[j].SiteComponentConfig.Name
	})
}
