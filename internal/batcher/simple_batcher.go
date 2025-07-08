package batcher

import "github.com/mach-composer/mach-composer-cli/internal/graph"

// simpleBatchFunc returns a BatchFunc that batches nodes based on their depth in the graph.
func simpleBatchFunc() BatchFunc {
	return func(g *graph.Graph) (map[int][]graph.Node, error) {
		batches := map[int][]graph.Node{}

		var sets = map[string][]graph.Path{}

		for _, n := range g.Vertices() {
			var route, _ = g.Routes(n.Path(), g.StartNode.Path())
			sets[n.Path()] = route
		}

		for k, routes := range sets {
			var mx int
			for _, route := range routes {
				if len(route) > mx {
					mx = len(route)
				}
			}
			n, _ := g.Vertex(k)
			batches[mx] = append(batches[mx], n)
		}

		return batches, nil
	}
}
