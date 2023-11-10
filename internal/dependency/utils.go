package dependency

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

func fetchSinks(nodeGraph NodeGraph) []Node {
	var sinks []Node

	am, _ := nodeGraph.AdjacencyMap()

	for k, v := range am {
		if len(v) == 0 {
			sink, _ := nodeGraph.Vertex(k)
			sinks = append(sinks, sink)
		}
	}

	return sinks
}

func fetchParents(nodeGraph NodeGraph, node Node) []Node {
	var parents []Node

	pm, _ := nodeGraph.PredecessorMap()

	for _, parent := range pm[node.Path()] {
		v, _ := nodeGraph.Vertex(parent.Source)
		parents = append(parents, v)
	}

	return parents
}
