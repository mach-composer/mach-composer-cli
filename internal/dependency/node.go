package dependency

import (
	"github.com/dominikbraun/graph"
)

type Node interface {
	Path() string
	Identifier() string
}

type nodeImpl struct {
	graph      graph.Graph[string, Node]
	path       string
	identifier string
}

func (n *nodeImpl) Path() string {
	return n.path
}

func (n *nodeImpl) Identifier() string {
	return n.identifier
}
