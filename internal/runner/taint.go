package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/mach-composer/mach-composer-cli/internal/hash"
	"github.com/rs/zerolog/log"
)

func determineTainted(oldHash string, n graph.Node, parentTainted bool) (bool, error) {
	// If a node has already been marked as tainted in a previous iteration, we don't need to check it again
	if n.Tainted() {
		return true, nil
	}

	// If a parent has been marked as tainted the current node is also tainted
	if parentTainted {
		return true, nil
	}

	h, err := n.Hash()
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to compute hash for %s", n.Path())
		return false, err
	}

	return h != oldHash, nil
}

func taintNode(ctx context.Context, hashFetcher hash.Handler, g *graph.Graph, path string, parentTainted bool) error {
	n, err := g.Vertex(path)
	if err != nil {
		return err
	}

	am, err := g.AdjacencyMap()
	if err != nil {
		return err
	}

	oldHash, err := hashFetcher.Fetch(ctx, n)
	if err != nil {
		return err
	}
	n.SetOldHash(oldHash)

	var isTainted = false
	if n.Type() != graph.ProjectType {
		isTainted, err = determineTainted(oldHash, n, parentTainted)
		if err != nil {
			return err
		}
	}
	n.SetTainted(isTainted)

	for _, child := range am[path] {
		if err = taintNode(ctx, hashFetcher, g, child.Target, isTainted); err != nil {
			return err
		}
	}

	return nil
}

func taintGraph(ctx context.Context, g *graph.Graph, hashFetcher hash.Handler) error {
	return taintNode(ctx, hashFetcher, g, g.StartNode.Path(), false)
}
