package runner

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/rs/zerolog/log"
)

type PlanOptions struct {
	Reuse      bool
	Lock       bool
	Components []string
	Site       string
}

func TerraformPlan(ctx context.Context, cfg *config.MachConfig, dg *dependency.Graph, opts *PlanOptions) error {
	if opts.Reuse {
		log.Warn().Msgf("Reuse option not implemented")
	}
	if len(opts.Components) > 0 {
		log.Warn().Msgf("Components option not implemented")
	}
	if opts.Site != "" {
		log.Warn().Msgf("Site option not implemented")
	}

	if err := terraformInitAll(ctx, dg); err != nil {
		return err
	}

	if err := batchRun(ctx, dg, cfg.MachComposer.Deployment.Runners,
		func(ctx context.Context, n dependency.Node, tfPath string) (string, error) {
			return terraformPlan(ctx, n, tfPath, opts)
		}); err != nil {
		return err
	}

	return nil
}

func terraformPlan(ctx context.Context, n dependency.Node, path string, opts *PlanOptions) (string, error) {
	cmd := []string{"plan"}

	if n.Type() == dependency.SiteComponentType {
		parents, err := n.Parents()
		if err != nil {
			return "", err
		}

		var skip = false
		var details []string
		for _, p := range parents {
			if p.Outputs().Type().HasAttribute("output") == false {
				details = append(details, fmt.Sprintf("parent %s has no output attribute", p.Path()))
				skip = true
			}
		}
		if skip == true {
			log.Warn().Strs("state", details).Msgf("Skipping plan for %s because of missing parent states", n.Path())
			return "", nil
		}
	}

	if opts.Lock == false {
		cmd = append(cmd, "-lock=false")
	}

	cmd = append(cmd, "-out=terraform.plan")
	return defaultRunTerraform(ctx, false, path, cmd...)
}
