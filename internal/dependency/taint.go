package dependency

import (
	"context"
)

func determineTainted(ctx context.Context, n Node, parentTainted bool) (bool, error) {
	if parentTainted {
		return true, nil
	}

	isTainted, err := n.HasConfigChanges(ctx)
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

func taintNode(ctx context.Context, g *Graph, path string, parentTainted bool) error {
	v, _ := g.Vertex(path)
	am, _ := g.AdjacencyMap()

	isTainted, err := determineTainted(ctx, v, parentTainted)
	if err != nil {
		return err
	}
	v.SetTainted(isTainted)

	for _, child := range am[path] {
		if err = taintNode(ctx, g, child.Target, isTainted); err != nil {
			return err
		}
	}

	return nil
}

func TaintNodes(ctx context.Context, g *Graph) error {
	return taintNode(ctx, g, g.StartNode.Path(), false)
}
