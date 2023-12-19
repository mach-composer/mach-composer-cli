package graph

func determineTainted(n Node, parentTainted bool) (bool, error) {
	// If a node has already been marked as tainted in a previous iteration, we don't need to check it again
	if n.Tainted() {
		return true, nil
	}

	// If a parent has been marked as tainted the current node is also tainted
	if parentTainted {
		return true, nil
	}

	// If a node has changes it is tainted
	return n.HasChanges()
}

func taintNode(g *Graph, path string, parentTainted bool) error {
	n, err := g.Vertex(path)
	if err != nil {
		return err
	}

	am, err := g.AdjacencyMap()
	if err != nil {
		return err
	}

	isTainted, err := determineTainted(n, parentTainted)
	if err != nil {
		return err
	}
	n.SetTainted(isTainted)

	for _, child := range am[path] {
		if err = taintNode(g, child.Target, isTainted); err != nil {
			return err
		}
	}

	return nil
}

// TaintNodes will mark all nodes as tainted that have changes or are dependent on a node with changes
func TaintNodes(g *Graph) error {
	return taintNode(g, g.StartNode.Path(), false)
}
