package runner

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

type PlanOptions struct {
	Reuse      bool
	Lock       bool
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

	if err := batchRun(ctx, dg, cfg.MachComposer.Deployment.Runners,
		func(ctx context.Context, n dependency.Node, tfPath string) (string, error) {
			return terraformPlan(ctx, n, tfPath)
		}); err != nil {
		return err
	}

	return nil
}

func terraformPlan(ctx context.Context, n dependency.Node, path string) (string, error) {
	log.Debug().Msgf("Running terraform plan for site %s", site.Identifier)

	if options.Reuse == false {
		if err := terraformInitSite(ctx, cfg, site, path); err != nil {
			return err
		}
	} else {
		log.Warn().Msgf("Skipping terraform init for site %s", site.Identifier)
	}

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

	if options.Lock == false {
		cmd = append(cmd, "-lock=false")
	}

	cmd = append(cmd, "-out=terraform.plan")
	return defaultRunTerraform(ctx, path, cmd...)
}
