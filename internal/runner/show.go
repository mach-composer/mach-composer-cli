package runner

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/rs/zerolog/log"

	"github.com/mach-composer/mach-composer-cli/internal/config"
)

type ShowPlanOptions struct {
	NoColor bool
	Site    string
}

func TerraformShow(ctx context.Context, cfg *config.MachConfig, dg *dependency.Graph, options *ShowPlanOptions) error {
	if err := batchRun(dg, dg.StartNode.Path(), func(n dependency.Node) error {
		tfPath := "deployments/" + n.Path()

		log.Info().Msgf("Showing %s", tfPath)

		return terraformShow(ctx, tfPath, options)
	}); err != nil {
		return err
	}

	return nil
}

func terraformShow(ctx context.Context, path string, options *ShowPlanOptions) error {
	filename, err := hasTerraformPlan(path)
	if err != nil {
		return err
	}
	if filename == "" {
		return fmt.Errorf("no plan found for path %s. Did you run `mach-composer plan`", path)
	}

	cmd := []string{"show", filename}
	if options.NoColor {
		cmd = append(cmd, "-no-color")
	}
	return runTerraform(ctx, path, cmd...)
}
