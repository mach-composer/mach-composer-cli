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
	return &SiteComponent{baseNode: baseNode{
		graph:          g,
		path:           path,
		identifier:     identifier,
		deploymentType: deploymentType,
		ancestor:       ancestor,
		typ:            SiteComponentType,
	}, SiteConfig: siteConfig, SiteComponentConfig: siteComponentConfig}
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
