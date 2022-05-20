package runner

import (
	"context"
	"github.com/labd/mach-composer/config"
	"github.com/sirupsen/logrus"
)

type ValidateOptions struct {
	Reuse bool
	Site  string
}

func TerraformValidate(cfg *config.MachConfig, locations map[string]string, options *ValidateOptions) {
	ctx := context.Background()

	for i := range cfg.Sites {
		site := cfg.Sites[i]

		if options.Site != "" && site.Identifier != options.Site {
			continue
		}

		TerraformValidateSite(ctx, cfg, &site, locations[site.Identifier], options)
	}
}

func TerraformValidateSite(ctx context.Context, cfg *config.MachConfig, site *config.Site, path string, options *ValidateOptions) {
	logrus.Debugf("Running terraform validate for site %s", site.Identifier)

	if !options.Reuse {
		RunTerraform(ctx, path, "init")
	}
	cmd := []string{"validate"}

	RunTerraform(ctx, path, cmd...)
}
