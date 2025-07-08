package graph

import (
	"github.com/dominikbraun/graph"
	"golang.org/x/exp/maps"
)

type Graph struct {
	graph.Graph[string, Node]
	StartNode Node
}

// Vertices returns all the vertex that are contained in the graph
func (g *Graph) Vertices() Vertices {
	var vertices Vertices

	m, _ := g.AdjacencyMap()

	keys := maps.Keys(m)

	for _, k := range keys {
		v, _ := g.Vertex(k)
		vertices = append(vertices, v)
	}

	return vertices
}

// Routes determines all the possible paths between two nodes
func (g *Graph) Routes(source, target string) ([]Path, error) {
	var routes []Path

	m, err := g.PredecessorMap()
	if err != nil {
		return routes, err
	}

	eg := m[source]

	for _, pathElement := range eg {
		p := []string{pathElement.Source}
		newRoutes := fetchPathsToTarget(pathElement.Source, target, m, p)
		routes = append(routes, newRoutes...)
	}

	return routes, nil
}

func (g *Graph) ExtractSubGraph(root Node) (*Graph, error) {
	// Create a new graph to hold the pruned subgraph
	newGraph := &Graph{
		Graph:     graph.New(func(n Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles()),
		StartNode: root,
	}
	if err := newGraph.AddVertex(root); err != nil {
		return nil, err
	}

	var addNodeAndChildren func(parent Node) error
	addNodeAndChildren = func(parent Node) error {
		children, err := parent.Children()
		if err != nil {
			return err
		}
		for _, child := range children {
			if err := newGraph.AddVertex(child); err != nil {
				return err
			}
			if err := newGraph.AddEdge(parent.Path(), child.Path()); err != nil {
				return err
			}
			if err := addNodeAndChildren(child); err != nil {
				return err
			}
		}
		return nil
	}

	if err := addNodeAndChildren(root); err != nil {
		return nil, err
	}

	return newGraph, nil
}

type Vertices []Node

func (v Vertices) Filter(t Type) Vertices {
	var nv Vertices

	for _, vx := range v {
		if vx.Type() == t {
			nv = append(nv, vx)
		}
	}
	return nv
}

func (v Vertices) FilterByIdentifier(identifier string) Node {
	for _, vx := range v {
		if vx.Identifier() == identifier {
			return vx
		}
	}

	return nil
}
