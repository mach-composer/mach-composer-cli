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

func TerraformProxy(ctx context.Context, cfg *config.MachConfig, dg *dependency.Graph, options *ProxyOptions) error {
	if err := batchRun(ctx, dg, dg.StartNode.Path(), cfg.MachComposer.Deployment.Runners,
		func(ctx context.Context, n dependency.Node) (string, error) {
			tfPath := "deployments/" + n.Path()

			log.Info().Msgf("Proxying command to %s", tfPath)

			return runTerraform(ctx, tfPath, options.Command...)
		}); err != nil {
		return err
	}

	return nil
}
