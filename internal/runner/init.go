package runner

import (
	"context"
	"github.com/rs/zerolog/log"

	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/lockfile"
)

type InitOptions struct {
	Site string
}

func TerraformInit(ctx context.Context, cfg *config.MachConfig, locations map[string]string, options *InitOptions) error {
	for i := range cfg.Sites {
		site := cfg.Sites[i]

		if options.Site != "" && site.Identifier != options.Site {
			continue
		}

		err := terraformInitSite(ctx, cfg, &site, locations[site.Identifier])
		if err != nil {
			return err
		}
	}
	return nil
}

func terraformInitSite(ctx context.Context, cfg *config.MachConfig, site *config.SiteConfig, path string) error {
	lf, err := lockfile.GetLock(cfg, path)
	if err != nil {
		return err
	}

	if !terraformIsInitialized(path) || lf.HasChanges(cfg) {
		log.Debug().Msgf("Running terraform init for site %s", site.Identifier)
		if err := defaultRunTerraform(ctx, path, "init"); err != nil {
			return err
		}
	}
	return nil
}
