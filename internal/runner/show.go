package runner

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type ShowPlanOptions struct {
	NoColor bool
}

func TerraformShow(ctx context.Context, cfg *config.MachConfig, dg *dependency.Graph, opt *ShowPlanOptions) error {
	if err := terraformInitAll(ctx, dg); err != nil {
		return err
	}

	if err := batchRun(ctx, dg, cfg.MachComposer.Deployment.Runners, func(ctx context.Context, n dependency.Node, tfPath string) (string, error) {
		return terraformShow(ctx, tfPath, opt)
	}); err != nil {
		return err
	}

	return nil
}

func terraformShow(ctx context.Context, path string, options *ShowPlanOptions) (string, error) {
	filename, err := hasTerraformPlan(path)
	if err != nil {
		return "", err
	}
	if filename == "" {
		return "", fmt.Errorf("no plan found for path %s. Did you run `mach-composer plan`", path)
	}

	cmd := []string{"show", filename}
	if options.NoColor {
		cmd = append(cmd, "-no-color")
	}
	return utils.RunTerraform(ctx, false, path, cmd...)
}
