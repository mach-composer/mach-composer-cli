package runner

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/lockfile"
)

type InitOptions struct {
	Site string
}

func TerraformInit(ctx context.Context, _ *config.MachConfig, dg *dependency.Graph, opt *InitOptions) error {
	if opt.Site != "" {
		log.Warn().Msgf("Site option not implemented")
	}
	return terraformInitAll(ctx, dg)
}

func terraformInitAll(ctx context.Context, g *dependency.Graph) error {
	var errChan = make(chan error, len(g.Vertices()))
	var wg = &sync.WaitGroup{}

	for _, n := range g.Vertices() {
		wg.Add(1)
		go func(n dependency.Node) {
			defer wg.Done()
			tfPath := "deployments/" + n.Path()
			hash, err := n.Hash()
			if err != nil {
				errChan <- err
				return
			}

			err = terraformInit(ctx, hash, tfPath)
			if err != nil {
				errChan <- err
				return
			}
		}(n)
	}
	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		var errors []error
		for err := range errChan {
			errors = append(errors, err)
		}

		return cli.NewGroupedError(
			fmt.Sprintf("failed initializing terraform projects (%d errors)", len(errors)), errors,
		)
	}

	return nil
}

func terraformInit(ctx context.Context, hash, path string) error {
	lf, err := lockfile.GetLock(hash, path)
	if err != nil {
		return err
	}

	if !terraformIsInitialized(path) || lf.HasChanges(hash) {
		if _, err = defaultRunTerraform(ctx, false, path, "init"); err != nil {
			return err
		}
	}
	return nil
}
