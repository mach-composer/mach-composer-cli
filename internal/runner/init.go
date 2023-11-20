package runner

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
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
	if err := batchRun(ctx, dg, cfg.MachComposer.Deployment.Runners, func(ctx context.Context, n dependency.Node, tfPath string) (string, error) {
		hash, err := n.Hash()
		if err != nil {
			return "", err
		}
		return terraformInit(ctx, hash, tfPath)
	}); err != nil {
		return err
	}

	return nil
}

func terraformInitAll(ctx context.Context, g *dependency.Graph) (string, error) {
	var out string
	for _, v := range g.Vertices() {
		tfPath := "deployments/" + v.Path()
		hash, err := v.Hash()
		if err != nil {
			return "", err
		}

		iOut, err := terraformInit(ctx, hash, tfPath)
		if err != nil {
			return "", err
		}
		out += fmt.Sprintf("%s\n", iOut)
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
		if out, err = utils.RunTerraform(ctx, path, "init"); err != nil {
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
