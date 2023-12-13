package dependency

import (
	"github.com/dominikbraun/graph"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
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
	return hashSiteComponentConfig(sc.SiteComponentConfig)
}

func (sc *SiteComponent) HasChanges() (bool, error) {
	hash, err := sc.Hash()
	if err != nil {
		return true, err
	}

	tfHash, err := utils.ParseHashOutput(sc.outputs)
	if err != nil {
		return false, err
	}

	return hash != tfHash, nil
}
