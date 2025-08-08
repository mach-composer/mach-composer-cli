package batcher

import (
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"golang.org/x/exp/maps"
)

// siteBatchFunc returns a BatchFunc that batches nodes based on their site order before considering their depth in
// the graph.
func siteBatchFunc(siteOrder []string) BatchFunc {
	return func(g *graph.Graph) (map[int][]graph.Node, error) {
		batches := map[int][]graph.Node{}

		var projects = g.Vertices().Filter(graph.ProjectType)
		if len(projects) != 1 {
			return nil, fmt.Errorf("expected 1 project, got %d", len(projects))
		}
		var project = projects[0]

		var sites = g.Vertices().Filter(graph.SiteType)

		batches[0] = []graph.Node{project}

		for _, siteIdentifier := range siteOrder {
			var sets = map[string][]graph.Path{}
			var site = sites.FilterByIdentifier(siteIdentifier)
			if site == nil {
				return nil, fmt.Errorf("site with identifier %s not found", siteIdentifier)
			}

			pg, err := g.ExtractSubGraph(site)
			if err != nil {
				return nil, err
			}

			for _, n := range pg.Vertices() {
				var route, _ = pg.Routes(n.Path(), site.Path())
				sets[n.Path()] = route
			}

			var siteBatches = map[int][]graph.Node{}

			for k, routes := range sets {
				var mx int
				for _, route := range routes {
					if len(route) > mx {
						mx = len(route)
					}
				}
				n, _ := pg.Vertex(k)
				siteBatches[mx] = append(siteBatches[mx], n)
			}

			// Get the highest int in the batches map
			var keys = maps.Keys(batches)
			var maxKey int
			for _, key := range keys {
				if key > maxKey {
					maxKey = key
				}
			}

			for k, v := range siteBatches {
				batches[maxKey+k+1] = append(batches[maxKey+k+1], v...)
			}
		}

		return batches, nil
	}
}
