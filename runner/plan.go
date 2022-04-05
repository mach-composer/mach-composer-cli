package runner

import (
	"context"
	"fmt"

	"github.com/labd/mach-composer/config"
	"github.com/sirupsen/logrus"
)

type PlanOptions struct {
	Reuse      bool
	Components []string
	Site       string
}

func TerraformPlan(cfg *config.MachConfig, locations map[string]string, options *PlanOptions) {
	ctx := context.Background()

	for i := range cfg.Sites {
		site := cfg.Sites[i]

		if options.Site != "" && site.Identifier != options.Site {
			continue
		}

		TerraformPlanSite(ctx, cfg, &site, locations[site.Identifier], options)
	}
}

func TerraformPlanSite(ctx context.Context, cfg *config.MachConfig, site *config.Site, path string, options *PlanOptions) {
	logrus.Debugf("Running terraform plan for site %s", site.Identifier)
	if !options.Reuse {
		RunTerraform(ctx, path, "init")
	}
	cmd := []string{"plan"}

	for _, component := range options.Components {
		cmd = append(cmd, fmt.Sprintf("-target=module.%s", component))
	}

	RunTerraform(ctx, path, cmd...)
}
