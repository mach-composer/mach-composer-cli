package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

type PlanOptions struct {
	Reuse      bool
	Components []string
	Site       string
}

func TerraformPlan(ctx context.Context, cfg *config.MachConfig, dg *dependency.Graph, opt *PlanOptions) error {
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

	if err := batchRun(ctx, dg, cfg.MachComposer.Deployment.Runners, func(ctx context.Context, n dependency.Node, tfPath string) (string, error) {
		return terraformPlan(ctx, tfPath)
	}); err != nil {
		return err
	}

	return nil
}

func terraformPlan(ctx context.Context, path string) (string, error) {
	cmd := []string{"plan"}

	cmd = append(cmd, "-out=terraform.plan")
	return utils.RunTerraform(ctx, false, path, cmd...)
}

func hasTerraformPlan(path string) (string, error) {
	filename := filepath.Join(path, "terraform.plan")
	if _, err := os.Stat(filename); err == nil {
		return filename, nil
	}
	return "", nil
}
