package dependency

type Batches map[int][]Node

type PathSet []Path

// Batch will batch the nodes for the given graph, so they can be run in parallel. A batch in this sense is a list of nodes
// that have no dependencies with each other. This currently happens with a very naive algorithm, where the batch number is
// determined by the longest Path from the root NodeImplementation to the NodeImplementation in question.
func Batch(g *Graph, s Node) (Batches, error) {
	batches := Batches{}

	var sets = map[string]PathSet{}

	for _, n := range g.Vertices() {
		var route, _ = g.Routes(n.Path(), s.Path())
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
