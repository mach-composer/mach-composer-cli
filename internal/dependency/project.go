package dependency

import (
	"github.com/dominikbraun/graph"
	"github.com/mach-composer/mach-composer-cli/internal/config"
)

type Project struct {
	baseNode
	ProjectConfig *config.MachConfig
}

func NewProject(g graph.Graph[string, Node], path, identifier string, deploymentType config.DeploymentType,
	projectConfig *config.MachConfig) *Project {
	return &Project{
		baseNode:      newBaseNode(g, path, identifier, ProjectType, nil, deploymentType),
		ProjectConfig: projectConfig,
	}
}

func (p *Project) Hash() (string, error) {
	return "", nil
}

// HasChanges always returns false, as any change in the project config will be picked up at the site or site component level
func (p *Project) HasChanges() (bool, error) {
	return false, nil
}
