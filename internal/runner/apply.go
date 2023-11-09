package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/rs/zerolog/log"

	"github.com/mach-composer/mach-composer-cli/internal/config"
)

type ApplyOptions struct {
	Reuse       bool
	Destroy     bool
	AutoApprove bool
	Components  []string
	Site        string
}

func TerraformApply(ctx context.Context, cfg *config.MachConfig, dg *dependency.Graph, options *ApplyOptions) error {
	if err := batchRun(dg, dg.StartNode.Path(), func(n dependency.Node) error {
		tfPath := "deployments/" + n.Path()

		log.Info().Msgf("Applying %s", tfPath)

		return terraformApply(ctx, cfg, tfPath, options)
	}); err != nil {
		return err
	}

	return nil
}

func terraformApply(ctx context.Context, cfg *config.MachConfig, path string, options *ApplyOptions) error {
	if err := terraformInit(ctx, cfg, path); err != nil {
		return err
	}

	cmd := []string{"apply"}

	if options.Destroy {
		cmd = append(cmd, "-destroy")
	}

	if options.AutoApprove {
		cmd = append(cmd, "-auto-approve")
	}

	// If there is a plan then we should use it.
	planFilename, err := hasTerraformPlan(path)
	if err != nil {
		return err
	}
	if planFilename != "" {
		cmd = append(cmd, planFilename)
	}

	return runTerraform(ctx, path, cmd...)
}
