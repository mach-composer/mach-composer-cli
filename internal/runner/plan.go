package runner

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"

	"github.com/mach-composer/mach-composer-cli/internal/config"
)

type PlanOptions struct {
	Reuse      bool
	Lock       bool
	Components []string
	Site       string
}

func TerraformPlan(ctx context.Context, cfg *config.MachConfig, locations map[string]string, options *PlanOptions) error {
	for i := range cfg.Sites {
		site := cfg.Sites[i]

		if options.Site != "" && site.Identifier != options.Site {
			continue
		}

		path := locations[site.Identifier]

		if err := terraformPlanSite(ctx, cfg, &site, path, options); err != nil {
			return err
		}
	}
	return nil
}

func terraformPlanSite(ctx context.Context, cfg *config.MachConfig, site *config.SiteConfig, path string, options *PlanOptions) error {
	log.Debug().Msgf("Running terraform plan for site %s", site.Identifier)

	if options.Reuse == false {
		if err := terraformInitSite(ctx, cfg, site, path); err != nil {
			return err
		}
	} else {
		log.Warn().Msgf("Skipping terraform init for site %s", site.Identifier)
	}

	cmd := []string{"plan"}
	for _, component := range options.Components {
		cmd = append(cmd, fmt.Sprintf("-target=module.%s", component))
	}

	if options.Lock == false {
		cmd = append(cmd, "-lock=false")
	}

	cmd = append(cmd, "-out=terraform.plan")
	return defaultRunTerraform(ctx, path, cmd...)
}
