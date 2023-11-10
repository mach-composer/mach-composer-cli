package runner

import (
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"golang.org/x/exp/maps"
	"sort"
)

type ExecutorFunc func(node dependency.Node) error

// batchRun will batch the nodes for the given graph, so they can be run in parallel. A batch in this sense is a list of nodes
// that have no dependencies with each other. This currently happens with a very naive algorithm, where the batch number is
// determined by the longest Path from the root Node to the Node in question.
func batchRun(g *dependency.Graph, start string, f ExecutorFunc) error {
	batches := map[int][]dependency.Node{}

	var sets = map[string][]dependency.Path{}

	for _, n := range g.Vertices() {
		var route, _ = g.Routes(n.Path(), start)
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

	keys := maps.Keys(batches)
	sort.Ints(keys)
	//TODO: add parallel running of batches based on number of workers
	for _, k := range keys[1:] {
		for _, n := range batches[k] {
			err := f(n)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
