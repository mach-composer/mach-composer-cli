package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/labd/mach-composer/internal/config"
)

type ApplyOptions struct {
	Reuse       bool
	Destroy     bool
	AutoApprove bool
	Components  []string
	Site        string
}

func TerraformApply(cfg *config.MachConfig, locations map[string]string, options *ApplyOptions) {
	ctx := context.Background()

	for i := range cfg.Sites {
		site := cfg.Sites[i]

		if options.Site != "" && site.Identifier != options.Site {
			continue
		}

		TerraformApplySite(ctx, cfg, &site, locations[site.Identifier], options)
	}
}

func TerraformProxy(cfg *config.MachConfig, locations map[string]string, siteName string, cmd []string) {
	ctx := context.Background()

	for i := range cfg.Sites {
		site := cfg.Sites[i]

		if siteName != "" && site.Identifier != siteName {
			continue
		}

		RunTerraform(ctx, locations[site.Identifier], cmd...)
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

	// If there is a plan then we should use it.
	if val := TerraformPlanDetect(path); val != "" {
		cmd = append(cmd, val)
	}

	RunTerraform(ctx, path, cmd...)
}

func TerraformPlanDetect(path string) string {
	filename := GeneratePlanName(path)
	filePath := filepath.Join(path, filename)

	if _, err := os.Stat(filePath); err == nil {
		return filename
	}
	return ""
}
