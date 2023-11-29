package runner

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
	"golang.org/x/sync/semaphore"
	"sort"
	"sync"
)

type ExecutorFunc func(ctx context.Context, node dependency.Node, tfPath string) (string, error)

// batchRun will batch the nodes for the given graph, so they can be run in parallel. A batch in this sense is a list of nodes
// that have no dependencies with each other. This currently happens with a very naive algorithm, where the batch number is
// determined by the longest Path from the root Node to the Node in question.
func batchRun(ctx context.Context, g *dependency.Graph, workers int, f ExecutorFunc) error {
	//	Load all the outputs for the nodes
	if err := g.LoadOutputs(ctx); err != nil {
		return err
	}

	// Check node hashes and parents and determine which nodes are tainted
	if err := dependency.TaintNodes(g); err != nil {
		return err
	}

	batches := map[int][]dependency.Node{}

	var sets = map[string][]dependency.Path{}

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

	keys := maps.Keys(batches)
	sort.Ints(keys)
	for i, k := range keys[1:] {
		log.Info().Msgf("Running batch %d with %d nodes", i, len(batches[k]))

		//Channel to collect errors
		errChan := make(chan error, len(batches[k]))

		//WaitGroup to make sure all go routines are finished
		wg := &sync.WaitGroup{}

		//Semaphore to limit the number of go routines
		sem := semaphore.NewWeighted(int64(workers))

		for _, n := range batches[k] {
			if err := sem.Acquire(ctx, 1); err != nil {
				return err
			}
			wg.Add(1)
			go func(ctx context.Context, n dependency.Node) {
				defer wg.Done()
				defer sem.Release(1)

				tfPath := "deployments/" + n.Path()

				if n.Tainted() == false {
					log.Info().Msgf("Skipping %s because it has no changes", tfPath)
					return
				}

				log.Info().Msgf("Running %s", tfPath)

				out, err := f(ctx, n, tfPath)
				//We allow the other elements of the batch to complete
				if err != nil {
					errChan <- err
					return
				}
				log.Info().Msg(out)
			}(ctx, n)
		}
		wg.Wait()
		close(errChan)

		if len(errChan) > 0 {
			var errors []error
			for err := range errChan {
				errors = append(errors, err)
			}

			return cli.NewGroupedError(fmt.Sprintf("batch run %d failed (%d errors)", i, len(errors)), errors)
		}

		log.Info().Msgf("Finished batch %d", i)
	}

	log.Info().Msgf("Finished all batches")

	return nil
}
