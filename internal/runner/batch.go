package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
	"golang.org/x/sync/semaphore"
	"sort"
	"sync"
)

type ExecutorFunc func(ctx context.Context, node dependency.Node) (string, error)

type BatchError struct {
	msg    string
	Errors []error
}

func NewBatchError(msg string, errors []error) *BatchError {
	return &BatchError{msg: msg, Errors: errors}
}

func (b *BatchError) Error() string {
	return b.msg
}

// batchRun will batch the nodes for the given graph, so they can be run in parallel. A batch in this sense is a list of nodes
// that have no dependencies with each other. This currently happens with a very naive algorithm, where the batch number is
// determined by the longest Path from the root Node to the Node in question.
func batchRun(ctx context.Context, g *dependency.Graph, start string, workers int, f ExecutorFunc) error {
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
	for i, k := range keys[1:] {
		log.Info().Msgf("Running batch %d with %d nodes", i, len(batches[k]))

		//Channel to collect execution results
		outChan := make(chan string, len(batches[k]))

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
			go func(n dependency.Node) {
				defer wg.Done()
				defer sem.Release(1)

				out, err := f(ctx, n)
				if err != nil {
					errChan <- err
				}
				outChan <- out
			}(n)
		}
		wg.Wait()
		close(errChan)
		close(outChan)

		if len(errChan) > 0 {
			var errors []error
			for err := range errChan {
				errors = append(errors, err)
			}

			return NewBatchError("batch run failed", errors)
		}

		for out := range outChan {
			log.Debug().Msg(out)
		}

		log.Info().Msgf("Finished batch %d", i)
	}

	return nil
}
