package runner

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/labd/mach-composer/internal/config"
)

type ShowOptions struct {
	Reuse      bool
	Components []string
	Site       string
}

func TerraformShow(ctx context.Context, cfg *config.MachConfig, locations map[string]string, options *PlanOptions) error {
	for i := range cfg.Sites {
		site := cfg.Sites[i]

		if options.Site != "" && site.Identifier != options.Site {
			continue
		}

		err := TerraformShowSite(ctx, cfg, &site, locations[site.Identifier], options)
		if err != nil {
			return err
		}
	}
	return nil
}

func TerraformShowSite(ctx context.Context, cfg *config.MachConfig, site *config.SiteConfig, path string, options *PlanOptions) error {
	filename, err := hasTerraformPlan(path)
	if err != nil {
		return err
	}
	if filename == "" {
		return fmt.Errorf("No plan found for site %s. Did you run `mach-composer plan`?", site.Identifier)
	}

	cmd := []string{"show", filename}
	log.Info().Msgf("Showing terraform plan for site %s", site.Identifier)
	return RunTerraform(ctx, path, cmd...)
}
