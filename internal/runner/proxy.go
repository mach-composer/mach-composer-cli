package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/rs/zerolog/log"
)

type ProxyOptions struct {
	Site    string
	Command []string
}

func TerraformProxy(ctx context.Context, cfg *config.MachConfig, dg *dependency.Graph, opt *ProxyOptions) error {
	if opt.Site != "" {
		log.Warn().Msgf("Site option not implemented")
	}

	if err := terraformInitAll(ctx, dg); err != nil {
		return err
	}

	if err := batchRun(ctx, dg, cfg.MachComposer.Deployment.Runners, func(ctx context.Context, n dependency.Node, tfPath string) (string, error) {
		return defaultRunTerraform(ctx, false, tfPath, opt.Command...)
	}); err != nil {
		return err
	}

	return nil
}
