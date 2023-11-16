package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mach-composer/mach-composer-cli/internal/config"
)

type PlanOptions struct {
	Reuse      bool
	Components []string
	Site       string
}

func TerraformPlan(ctx context.Context, cfg *config.MachConfig, dg *dependency.Graph, _ *PlanOptions) error {
	if err := batchRun(ctx, dg, dg.StartNode.Path(), cfg.MachComposer.Deployment.Runners,
		func(ctx context.Context, n dependency.Node) (string, error) {
			tfPath := "deployments/" + n.Path()

			log.Info().Msgf("Planning %s", tfPath)

			return terraformPlan(ctx, cfg, tfPath)
		}); err != nil {
		return err
	}

	return nil
}

func terraformPlan(ctx context.Context, cfg *config.MachConfig, path string) (string, error) {
	//TODO: deal with situation where a terraform file refers to a remote state that does not already exist
	iOut, err := terraformInit(ctx, cfg, path)
	if err != nil {
		return "", err
	}

	cmd := []string{"plan"}

	cmd = append(cmd, "-out=terraform.plan")
	rOut, err := runTerraform(ctx, path, cmd...)
	if err != nil {
		return "", err
	}

	return strings.Join([]string{iOut, rOut}, "\n"), nil
}

func hasTerraformPlan(path string) (string, error) {
	filename := filepath.Join(path, "terraform.plan")
	if _, err := os.Stat(filename); err == nil {
		return filename, nil
	}
	return "", nil
}
