package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"strings"
)

type ApplyOptions struct {
	Destroy     bool
	AutoApprove bool
	Workers     int
}

func TerraformApply(ctx context.Context, dg *dependency.Graph, opts *ApplyOptions) error {
	if err := batchRun(ctx, dg, opts.Workers, func(ctx context.Context, n dependency.Node) (string, error) {
		return terraformApply(ctx, n.Path(), opts.Destroy, opts.AutoApprove)
	}); err != nil {
		return err
	}

	return nil
}

func terraformApply(ctx context.Context, path string, destroy, autoApprove bool) (string, error) {
	cmd := []string{"apply"}

	if destroy {
		cmd = append(cmd, "-destroy")
	}

	if autoApprove {
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
