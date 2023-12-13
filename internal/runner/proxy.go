package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type ProxyOptions struct {
	Command []string
	Workers int
}

func TerraformProxy(ctx context.Context, dg *dependency.Graph, opts *ProxyOptions) error {
	if err := batchRun(ctx, dg, opts.Workers, func(ctx context.Context, n dependency.Node) (string, error) {
		return utils.RunTerraform(ctx, false, n.Path(), opts.Command...)
	}); err != nil {
		return err
	}

	return nil
}
