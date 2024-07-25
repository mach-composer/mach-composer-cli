package graph

import (
	"github.com/dominikbraun/graph"
)

type Path []string

func fetchPathsToTarget(source, target string, pm map[string]map[string]graph.Edge[string], currentPath Path) []Path {
	var paths []Path
	parents := pm[source]
	if len(parents) == 0 {
		return []Path{currentPath}
	}

	currentPath = append(currentPath, source)

	for _, parent := range parents {
		if parent.Source == target {
			paths = []Path{currentPath}
		}
		newPaths := fetchPathsToTarget(parent.Source, target, pm, currentPath)
		paths = append(paths, newPaths...)
	}

	return paths
}

func pruneBranch(g *Graph, n Node) error {
	var pErr error
	graph.DFS(g.Graph, n.Path(), func(p string) bool {
		n, err := g.Graph.Vertex(p)
		if err != nil {
			pErr = err
			return true
		}

		e, err := g.Edges()
		if err != nil {
			pErr = err
			return true
		}

		var edges []graph.Edge[string]

		for _, edge := range e {
			if edge.Target == n.Path() {
				edges = append(edges, edge)
			}

			if edge.Source == n.Path() {
				edges = append(edges, edge)
			}
		}

		for _, edge := range edges {
			err = g.RemoveEdge(edge.Source, edge.Target)
			if err != nil {
				pErr = err
				return true
			}
		}

		err = g.RemoveVertex(n.Path())
		if err != nil {
			pErr = err
			return true
		}

		return false
	})

	return pErr
}
