package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type ProxyOptions struct {
	Command []string
}

func TerraformProxy(ctx context.Context, cfg *config.MachConfig, dg *dependency.Graph, opt *ProxyOptions) error {
	if err := terraformInitAll(ctx, dg); err != nil {
		return err
	}

	if err := batchRun(ctx, dg, cfg.MachComposer.Deployment.Runners, func(ctx context.Context, n dependency.Node, tfPath string) (string, error) {
		return utils.RunTerraform(ctx, false, tfPath, opt.Command...)
	}); err != nil {
		return err
	}

	return nil
}
