package tree

import (
	"github.com/dominikbraun/graph"
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
