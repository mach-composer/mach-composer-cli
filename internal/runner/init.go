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

func TerraformInit(ctx context.Context, _ *config.MachConfig, dg *dependency.Graph, _ *InitOptions) error {
	out, err := terraformInitAll(ctx, dg)
	if err != nil {
		return err
	}
	log.Info().Msg(out)
	return nil
}

func terraformInitAll(ctx context.Context, g *dependency.Graph) (string, error) {
	var out string
	var errChan = make(chan error, len(g.Vertices()))
	var wg = &sync.WaitGroup{}
	var mu = &sync.Mutex{}

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

			iOut, err := terraformInit(ctx, hash, tfPath)
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			out += fmt.Sprintf("%s\n", iOut)
			mu.Unlock()
		}(n)
	}
	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		var errors []error
		for err := range errChan {
			errors = append(errors, err)
		}

		return "", cli.NewGroupedError(
			fmt.Sprintf("failed initializing terraform projects (%d errors)", len(errors)), errors,
		)
	}

	return out, nil
}

func terraformInit(ctx context.Context, hash, path string) (string, error) {
	lf, err := lockfile.GetLock(hash, path)
	if err != nil {
		return "", err
	}

	var out string
	if !terraformIsInitialized(path) || lf.HasChanges(hash) {
		if out, err = utils.RunTerraform(ctx, false, path, "init"); err != nil {
			return "", err
		}
	}
	return out, nil
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
