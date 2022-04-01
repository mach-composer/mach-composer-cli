package runner

import (
	"context"
	"fmt"

	"github.com/labd/mach-composer-go/config"
)

type PlanOptions struct {
	Reuse      bool
	Components []string
}

func TerraformPlan(cfg *config.MachConfig, locations map[string]string, options *PlanOptions) {
	ctx := context.Background()

	for i := range cfg.Sites {
		site := cfg.Sites[i]
		TerraformPlanSite(ctx, cfg, &site, locations[site.Identifier], options)
	}
}

func TerraformPlanSite(ctx context.Context, cfg *config.MachConfig, site *config.Site, path string, options *PlanOptions) {
	if !options.Reuse {
		RunTerraform(ctx, path, "init")
	}
	cmd := []string{"plan"}

	for _, component := range options.Components {
		cmd = append(cmd, fmt.Sprintf("-target=module.%s", component))
	}

	RunTerraform(ctx, path, cmd...)
}
