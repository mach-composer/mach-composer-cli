package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/lockfile"
)

type InitOptions struct {
	Site string
}

func TerraformInit(ctx context.Context, cfg *config.MachConfig, dg *dependency.Graph, _ *InitOptions) error {
	if err := batchRun(ctx, dg, dg.StartNode.Path(), cfg.MachComposer.Deployment.Runners,
		func(ctx context.Context, n dependency.Node) (string, error) {
			tfPath := "deployments/" + n.Path()

			log.Info().Msgf("Initializing %s", tfPath)

			return terraformInit(ctx, cfg, tfPath)
		}); err != nil {
		return err
	}

	return nil
}

func terraformInit(ctx context.Context, cfg *config.MachConfig, path string) (string, error) {
	lf, err := lockfile.GetLock(cfg.ConfigHash, path)
	if err != nil {
		return "", err
	}

	var out string
	if !terraformIsInitialized(path) || lf.HasChanges(cfg.ConfigHash) {
		if out, err = runTerraform(ctx, path, "init"); err != nil {
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
