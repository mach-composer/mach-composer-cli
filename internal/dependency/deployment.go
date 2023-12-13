package dependency

import (
	"errors"
	"fmt"
	"github.com/dominikbraun/graph"
	"github.com/mach-composer/mach-composer-cli/internal/config"
)

// ToDeploymentGraph converts a MachConfig to a Graph ready for deployment.
// This means that all nodes that are not independently deployable are pruned from the graph.
func ToDeploymentGraph(cfg *config.MachConfig, outPath string) (*Graph, error) {
	g, err := ToDependencyGraph(cfg, outPath)
	if err != nil {
		return nil, err
	}

	if err := validateDeployment(g); err != nil {
		return nil, err
	}

	// Remove all nodes that are not independent to site node
	if err = reduceNodes(g); err != nil {
		return nil, err
	}

	return g, nil
}

func reduceNodes(g *Graph) error {
	var pErr error
	if err := graph.BFS(g.Graph, g.StartNode.Path(), func(p string) bool {
		n, _ := g.Graph.Vertex(p)

		if !n.Independent() {
			siteNode, ok := n.Ancestor().(*Site)
			if !ok {
				pErr = fmt.Errorf("node %s is expected to have site as parent", n.Path())
				return true
			}

			siteComponentNode, ok := n.(*SiteComponent)
			if !ok {
				pErr = fmt.Errorf("node %s is expected to be a site component", n.Path())
				return true
			}

			siteNode.NestedSiteComponentConfigs = append(
				siteNode.NestedSiteComponentConfigs,
				siteComponentNode.SiteComponentConfig,
			)

			am, _ := g.Graph.AdjacencyMap()
			pm, _ := g.Graph.PredecessorMap()

			childEdges := am[p]
			parentEdges := pm[p]

			for _, childEdge := range childEdges {
				if err := g.Graph.RemoveEdge(childEdge.Source, childEdge.Target); err != nil {
					pErr = err
					return false
				}
			}

			for _, parentEdge := range parentEdges {
				if err := g.Graph.RemoveEdge(parentEdge.Source, parentEdge.Target); err != nil {
					pErr = err
					return false
				}
			}

			for _, childEdge := range childEdges {
				for _, parentEdge := range parentEdges {

					_, err := g.Graph.Edge(parentEdge.Source, childEdge.Target)
					if err != nil && !errors.Is(err, graph.ErrEdgeNotFound) {
						pErr = err
						return false
					}

					if err != nil && errors.Is(err, graph.ErrEdgeNotFound) {
						if err := g.Graph.AddEdge(parentEdge.Source, childEdge.Target); err != nil {
							pErr = err
							return false
						}
					}
				}
			}

			if err := g.Graph.RemoveVertex(n.Path()); err != nil {
				pErr = err
				return true
			}
		}

		return false
	}); err != nil {
		return err
	}

	for _, v := range g.Vertices() {
		v.resetGraph(g.Graph)
	}

	return pErr
}
