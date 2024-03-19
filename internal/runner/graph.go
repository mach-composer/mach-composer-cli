package runner

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/batcher"
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/mach-composer/mach-composer-cli/internal/hash"
	"github.com/mach-composer/mach-composer-cli/internal/terraform"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
	"golang.org/x/sync/semaphore"
	"sort"
	"sync"
)

type (
	//executorFunc is a function that executes an arbitrary command on a node
	executorFunc func(ctx context.Context, node graph.Node) (string, error)
)

// GraphRunner will run a set of commands on a graph of nodes. Untainted nodes (no changes) will be skipped.
// The nodes are batched based on a batching function, and all nodes in the same batch will be run in parallel.
type GraphRunner struct {
	workers int
	batch   batcher.BatchFunc
	hash    hash.Handler
}

func NewGraphRunner(batcher batcher.BatchFunc, hashHandler hash.Handler, workers int) *GraphRunner {
	return &GraphRunner{
		workers: workers,
		batch:   batcher,
		hash:    hashHandler,
	}
}

func (gr *GraphRunner) run(ctx context.Context, g *graph.Graph, f executorFunc, ignoreChangeDetection bool) error {
	if err := taintGraph(ctx, g, gr.hash); err != nil {
		return err
	}

	batches := gr.batch(g)

	keys := maps.Keys(batches)
	sort.Ints(keys)
	for i, k := range keys[1:] {
		log.Info().Msgf("Running batch %d with %d nodes", i, len(batches[k]))

		errChan := make(chan error, len(batches[k]))
		wg := &sync.WaitGroup{}
		sem := semaphore.NewWeighted(int64(gr.workers))

		for _, n := range batches[k] {
			if n.Tainted() == false && ignoreChangeDetection == false {
				log.Info().Msgf("Skipping %s because it has no changes", n.Identifier())
				continue
			}

			if err := sem.Acquire(ctx, 1); err != nil {
				return err
			}
			wg.Add(1)
			go func(ctx context.Context, n graph.Node) {
				defer wg.Done()
				defer sem.Release(1)

				log.Info().Msgf("Running command on %s", n.Identifier())

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

func (gr *GraphRunner) TerraformApply(ctx context.Context, dg *graph.Graph, opts *ApplyOptions) error {
	if err := gr.run(ctx, dg, func(ctx context.Context, n graph.Node) (string, error) {
		if !terraformIsInitialized(n.Path()) || opts.ForceInit {
			log.Info().Msgf("Running terraform init for %s", n.Path())
			if out, err := terraform.Init(ctx, n.Path()); err != nil {
				return out, err
			}
		} else {
			log.Info().Msgf("Skipping terraform init for %s", n.Path())
		}

		out, err := terraform.Apply(ctx, n.Path(), opts.Destroy, opts.AutoApprove)
		if err != nil {
			return out, err
		}

		log.Info().Msgf("Storing new hash for %s", n.Path())
		if err = gr.hash.Store(ctx, n); err != nil {
			log.Warn().Err(err).Msgf("Failed to store hash for %s", n.Identifier())
		}
		return out, nil

	}, opts.IgnoreChangeDetection); err != nil {
		return err
	}

	return nil
}

func (gr *GraphRunner) TerraformPlan(ctx context.Context, dg *graph.Graph, opts *PlanOptions) error {
	if err := gr.run(ctx, dg, func(ctx context.Context, n graph.Node) (string, error) {
		if !terraformIsInitialized(n.Path()) || opts.ForceInit {
			log.Info().Msgf("Running terraform init for %s", n.Path())
			if out, err := terraform.Init(ctx, n.Path()); err != nil {
				return out, err
			}
		} else {
			log.Info().Msgf("Skipping terraform init for %s", n.Path())
		}

		canPlan, err := terraformCanPlan(ctx, n)
		if err != nil {
			return "", err
		}

		if !canPlan {
			log.Info().Msgf("Skipping planning %s because it has missing outputs", n.Path())
			return "", nil
		}

		return terraform.Plan(ctx, n.Path(), opts.Lock)
	}, opts.IgnoreChangeDetection); err != nil {
		return err
	}

	return nil
}

func (gr *GraphRunner) TerraformProxy(ctx context.Context, dg *graph.Graph, opts *ProxyOptions) error {
	if err := gr.run(ctx, dg, func(ctx context.Context, n graph.Node) (string, error) {
		if !terraformIsInitialized(n.Path()) {
			return "", fmt.Errorf("terraform is not initialized for %s. Please run init beforehand", n.Path())
		}

		return utils.RunTerraform(ctx, n.Path(), false, opts.Command...)
	}, opts.IgnoreChangeDetection); err != nil {
		return err
	}

	return nil
}

func (gr *GraphRunner) TerraformShow(ctx context.Context, dg *graph.Graph, opts *ShowPlanOptions) error {
	if err := gr.run(ctx, dg, func(ctx context.Context, n graph.Node) (string, error) {
		return terraform.Show(ctx, n.Path(), opts.NoColor)
	}, opts.IgnoreChangeDetection); err != nil {
		return err
	}

	return nil
}

func (gr *GraphRunner) TerraformInit(ctx context.Context, dg *graph.Graph) error {
	if err := gr.run(ctx, dg, func(ctx context.Context, n graph.Node) (string, error) {
		return terraform.Init(ctx, n.Path())
	}, true); err != nil {
		return err
	}

	return nil
}
