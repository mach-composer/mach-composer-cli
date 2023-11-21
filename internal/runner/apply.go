package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/rs/zerolog/log"
	"strings"
)

type ApplyOptions struct {
	Reuse       bool
	Destroy     bool
	AutoApprove bool
	Components  []string
	Site        string
}

func TerraformApply(ctx context.Context, cfg *config.MachConfig, dg *dependency.Graph, opt *ApplyOptions) error {
	if opt.Reuse {
		log.Warn().Msgf("Reuse option not implemented")
	}
	if len(opt.Components) > 0 {
		log.Warn().Msgf("Components option not implemented")
	}
	if opt.Site != "" {
		log.Warn().Msgf("Site option not implemented")
	}

	if err := terraformInitAll(ctx, dg); err != nil {
		return err
	}

	if err := batchRun(ctx, dg, cfg.MachComposer.Deployment.Runners, func(ctx context.Context, _ dependency.Node,
		tfPath string) (string, error) {
		return terraformApply(ctx, tfPath, opt)
	}); err != nil {
		return err
	}

	return nil
}

func terraformApply(ctx context.Context, path string, opt *ApplyOptions) (string, error) {
	cmd := []string{"apply"}

	if opt.Destroy {
		cmd = append(cmd, "-destroy")
	}

	if opt.AutoApprove {
		cmd = append(cmd, "-auto-approve")
	}

	// If there is a plan then we should use it.
	planFilename, err := hasTerraformPlan(path)
	if err != nil {
		return "", err
	}
	if planFilename != "" {
		cmd = append(cmd, strings.TrimPrefix(planFilename, path+"/"))
	}

	return utils.RunTerraform(ctx, false, path, cmd...)
}
