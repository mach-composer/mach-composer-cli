package runner

import (
	"context"
	"fmt"

	"github.com/labd/mach-composer-go/config"
)

type ApplyOptions struct {
	Reuse       bool
	Destroy     bool
	AutoApprove bool
	Components  []string
}

func TerraformApply(cfg *config.MachConfig, locations map[string]string, options *ApplyOptions) {
	ctx := context.Background()

	for i := range cfg.Sites {
		site := cfg.Sites[i]
		TerraformApplySite(ctx, cfg, &site, locations[site.Identifier], options)
	}
}

func TerraformApplySite(ctx context.Context, cfg *config.MachConfig, site *config.Site, path string, options *ApplyOptions) {

	if !options.Reuse {
		RunTerraform(ctx, path, "init")
	}

	cmd := []string{"apply"}

	if options.Destroy {
		cmd = append(cmd, "-destroy")
	}

	if options.AutoApprove {
		cmd = append(cmd, "-auto-approve")
	}

	for _, component := range options.Components {
		cmd = append(cmd, fmt.Sprintf("-target=module.%s", component))
	}

	RunTerraform(ctx, path, cmd...)
}
