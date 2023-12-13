package runner

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
	"golang.org/x/sync/semaphore"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

const PlanFile = "terraform.plan"

// TODO: write tests
func batchNodes(g *dependency.Graph) map[int][]dependency.Node {
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

	return batches
}

// LoadOutputs loads the outputs for all nodes in the graph in parallel
func loadOutputs(ctx context.Context, g *dependency.Graph) error {
	wg := &sync.WaitGroup{}
	errChan := make(chan error, len(g.Vertices()))

	for _, n := range g.Vertices() {
		wg.Add(1)

		go func(ctx context.Context, n dependency.Node) {
			defer wg.Done()
			val, err := utils.GetTerraformOutputs(ctx, n.Path())
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

type ExecutorFunc func(ctx context.Context, node dependency.Node) (string, error)

// batchRun will batch the nodes for the given graph, so they can be run in parallel. A batch in this sense is a list of nodes
// that have no dependencies with each other. This currently happens with a very naive algorithm, where the batch number is
// determined by the longest Path from the root Node to the Node in question.
// TODO: write tests
func batchRun(ctx context.Context, g *dependency.Graph, workers int, f ExecutorFunc) error {
	if err := loadOutputs(ctx, g); err != nil {
		return err
	}

	if err := dependency.TaintNodes(g); err != nil {
		return err
	}

	batches := batchNodes(g)

	keys := maps.Keys(batches)
	sort.Ints(keys)
	for i, k := range keys[1:] {
		log.Info().Msgf("Running batch %d with %d nodes", i, len(batches[k]))

		errChan := make(chan error, len(batches[k]))
		wg := &sync.WaitGroup{}
		sem := semaphore.NewWeighted(int64(workers))

		for _, n := range batches[k] {
			if err := sem.Acquire(ctx, 1); err != nil {
				return err
			}
			wg.Add(1)
			go func(ctx context.Context, n dependency.Node) {
				defer wg.Done()
				defer sem.Release(1)

				if n.Tainted() == false {
					log.Info().Msgf("Skipping %s because it has no changes", n.Identifier())
					return
				}

				log.Info().Msgf("Running apply for %s", n.Identifier())

				out, err := f(ctx, n)
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

func hasTerraformPlan(path string) (string, error) {
	filename := filepath.Join(path, PlanFile)
	if _, err := os.Stat(filename); err == nil {
		return filename, nil
	}
	return "", nil
}

func terraformIsInitialized(path string) bool {
	tfLockFile := filepath.Join(path, ".terraform.lock.hcl")
	if _, err := os.Stat(tfLockFile); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatal().Err(err)
	}
	return true
}
