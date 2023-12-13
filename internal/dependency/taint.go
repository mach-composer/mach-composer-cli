package dependency

func determineTainted(n Node, parentTainted bool) (bool, error) {
	if parentTainted {
		return true, nil
	}

	isTainted, err := n.HasChanges()
	if err != nil {
		return false, err
	}

	if isTainted {
		return true, nil
	}

	if n.Tainted() {
		return true, nil
	}

	return false, nil
}

func taintNode(g *Graph, path string, parentTainted bool) error {
	v, _ := g.Vertex(path)
	am, _ := g.AdjacencyMap()

	isTainted, err := determineTainted(v, parentTainted)
	if err != nil {
		return err
	}
	v.SetTainted(isTainted)

	for _, child := range am[path] {
		if err = taintNode(g, child.Target, isTainted); err != nil {
			return err
		}
	}

	return nil
}

// TaintNodes will mark all nodes as tainted that have changes or are dependent on a node with changes
// TODO: write tests
func TaintNodes(g *Graph) error {
	return taintNode(g, g.StartNode.Path(), false)
}
