package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/rs/zerolog/log"
)

type ApplyOptions struct {
	Reuse       bool
	Destroy     bool
	AutoApprove bool
	Components  []string
	Site        string
}

func TerraformApply(ctx context.Context, cfg *config.MachConfig, dg *dependency.Graph, options *ApplyOptions) error {
	out, err := terraformInitAll(ctx, dg)
	if err != nil {
		return err
	}
	log.Debug().Msg(out)

	if err = batchRun(ctx, dg, cfg.MachComposer.Deployment.Runners, func(ctx context.Context, _ dependency.Node, tfPath string) (string, error) {
		return terraformApply(ctx, tfPath, options)
	}); err != nil {
		return err
	}

	return nil
}

func terraformApply(ctx context.Context, path string, options *ApplyOptions) (string, error) {
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
		return "", err
	}
	if planFilename != "" {
		cmd = append(cmd, planFilename)
	}

	return utils.RunTerraform(ctx, path, cmd...)
}
