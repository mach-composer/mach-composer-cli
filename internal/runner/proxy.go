package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/rs/zerolog/log"
)

type ProxyOptions struct {
	Site    string
	Command []string
}

func TerraformProxy(ctx context.Context, cfg *config.MachConfig, dg *dependency.Graph, options *ProxyOptions) error {
	out, err := terraformInitAll(ctx, dg)
	if err != nil {
		return err
	}
	log.Debug().Msg(out)

	if err := batchRun(ctx, dg, cfg.MachComposer.Deployment.Runners, func(ctx context.Context, n dependency.Node, tfPath string) (string, error) {
		return utils.RunTerraform(ctx, tfPath, options.Command...)
	}); err != nil {
		return err
	}

	return nil
}
