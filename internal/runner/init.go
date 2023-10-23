package runner

import (
	"context"
	"os"
	"path/filepath"

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
	lockfile, err := lockfile.GetLock(cfg, path)
	if err != nil {
		return err
	}

	if !terraformIsInitialized(path) || lockfile.HasChanges(cfg) {
		log.Debug().Msgf("Running terraform init for site %s", site.Identifier)
		if err := RunTerraform(ctx, path, "init"); err != nil {
			return err
		}
	}
	return nil
}

func terraformIsInitialized(path string) bool {
	tfLockFile := filepath.Join(path, ".terraform.lock.hcl")
	if _, err := os.Stat(tfLockFile); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatal().Err(err)
	}
	return true
}
