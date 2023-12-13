package runner

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type ShowPlanOptions struct {
	NoColor bool
	Workers int
}

func TerraformShow(ctx context.Context, dg *dependency.Graph, opts *ShowPlanOptions) error {
	if err := batchRun(ctx, dg, opts.Workers, func(ctx context.Context, n dependency.Node) (string, error) {
		return terraformShow(ctx, n.Path(), opts.NoColor)
	}); err != nil {
		return err
	}

	return nil
}

func terraformShow(ctx context.Context, path string, noColor bool) (string, error) {
	filename, err := hasTerraformPlan(path)
	if err != nil {
		return "", err
	}
	if filename == "" {
		return "", fmt.Errorf("no plan found for path %s. Did you run `mach-composer plan`", path)
	}

	cmd := []string{"show", filename}
	if noColor {
		cmd = append(cmd, "-no-color")
	}
	return utils.RunTerraform(ctx, false, path, cmd...)
}
