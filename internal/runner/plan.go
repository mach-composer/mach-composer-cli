package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/mach-composer/mach-composer-cli/internal/config"
)

type PlanOptions struct {
	Reuse      bool
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

		if err := TerraformPlanSite(ctx, cfg, &site, path, options); err != nil {
			return err
		}
	}
	return nil
}

func TerraformPlanSite(ctx context.Context, cfg *config.MachConfig, site *config.SiteConfig, path string, options *PlanOptions) error {
	log.Debug().Msgf("Running terraform plan for site %s", site.Identifier)

	if err := terraformInitSite(ctx, cfg, site, path); err != nil {
		return err
	}

	cmd := []string{"plan"}
	for _, component := range options.Components {
		cmd = append(cmd, fmt.Sprintf("-target=module.%s", component))
	}

	cmd = append(cmd, "-out=terraform.plan")
	return RunTerraform(ctx, path, cmd...)
}

func hasTerraformPlan(path string) (string, error) {
	filename := filepath.Join(path, "terraform.plan")
	if _, err := os.Stat(filename); err == nil {
		return filename, nil
	}
	return "", nil
}
