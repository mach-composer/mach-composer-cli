package graph

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/zclconf/go-cty/cty"
	"sync"
)

type (
	outputLoader func(ctx context.Context, path string) (cty.Value, error)
)

// LoadOutputs loads the outputs for all nodes in the graph in parallel
func LoadOutputs(ctx context.Context, g *Graph, loader outputLoader) error {
	wg := &sync.WaitGroup{}
	errChan := make(chan error, len(g.Vertices()))

	for _, n := range g.Vertices() {
		wg.Add(1)

		go func(ctx context.Context, n Node) {
			defer wg.Done()
			val, err := loader(ctx, n.Path())
			if err != nil {
				errChan <- err
				return
			}
			n.SetOutputs(val)
		}(ctx, n)
	}
	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		var errors []error
		for err := range errChan {
			errors = append(errors, err)
		}

		return cli.NewGroupedError(fmt.Sprintf("failed loading outputs (%d errors)", len(errors)), errors)
	}

	return nil
}
