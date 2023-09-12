package dependency

import (
	"github.com/dominikbraun/graph"
	"github.com/mach-composer/mach-composer-cli/internal/config"
)

type Node interface {
	Path() string
}

type nodeImpl struct {
	graph graph.Graph[string, Node]
	path  string
}

func (n *nodeImpl) Path() string {
	return n.path
}

type Project struct {
	nodeImpl
	Config *config.MachConfig
}

type Site struct {
	nodeImpl
	Config *config.SiteConfig
}

type SiteComponent struct {
	nodeImpl
	Config *config.SiteComponentConfig
}
