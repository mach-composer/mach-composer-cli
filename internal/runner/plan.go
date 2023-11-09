package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"

	"github.com/mach-composer/mach-composer-cli/internal/config"
)

type PlanOptions struct {
	Reuse      bool
	Components []string
	Site       string
}

func TerraformPlan(ctx context.Context, cfg *config.MachConfig, dg *dependency.Graph, options *PlanOptions) error {
	if err := batchRun(dg, dg.StartNode.Path(), func(n dependency.Node) error {
		tfPath := "deployments/" + n.Path()

		log.Info().Msgf("Planning %s", tfPath)

		return terraformPlan(ctx, cfg, tfPath)
	}); err != nil {
		return err
	}

	return nil
}

func terraformPlan(ctx context.Context, cfg *config.MachConfig, path string) error {
	//TODO: deal with situation where a terraform file refers to a remote state that does not already exist
	if err := terraformInit(ctx, cfg, path); err != nil {
		return err
	}

	cmd := []string{"plan"}

	cmd = append(cmd, "-out=terraform.plan")
	return runTerraform(ctx, path, cmd...)
}

func hasTerraformPlan(path string) (string, error) {
	filename := filepath.Join(path, "terraform.plan")
	if _, err := os.Stat(filename); err == nil {
		return filename, nil
	}
	return "", nil
}
