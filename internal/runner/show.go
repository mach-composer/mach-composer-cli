package runner

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/mach-composer/mach-composer-cli/internal/config"
)

type ShowPlanOptions struct {
	NoColor bool
	Site    string
}

func TerraformShow(ctx context.Context, cfg *config.MachConfig, locations map[string]string, options *ShowPlanOptions) error {
	for i := range cfg.Sites {
		site := cfg.Sites[i]

		if options.Site != "" && site.Identifier != options.Site {
			continue
		}

		err := terraformShowSite(ctx, cfg, &site, locations[site.Identifier], options)
		if err != nil {
			return err
		}
	}
	return nil
}

func terraformShowSite(ctx context.Context, cfg *config.MachConfig, site *config.SiteConfig, path string, options *ShowPlanOptions) error {
	filename, err := hasTerraformPlan(path)
	if err != nil {
		return err
	}
	if filename == "" {
		return fmt.Errorf("No plan found for site %s. Did you run `mach-composer plan`?", site.Identifier)
	}

	cmd := []string{"show", filename}
	if options.NoColor {
		cmd = append(cmd, "-no-color")
	}
	log.Info().Msgf("Showing terraform plan for site %s", site.Identifier)
	return RunTerraform(ctx, path, cmd...)
}
