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

func TerraformApply(ctx context.Context, cfg *config.MachConfig, locations map[string]string, options *ApplyOptions) error {
	for i := range cfg.Sites {
		site := cfg.Sites[i]

		if options.Site != "" && site.Identifier != options.Site {
			continue
		}

		if err := TerraformApplySite(ctx, cfg, &site, locations[site.Identifier], options); err != nil {
			return err
		}
	}
	return nil
}

type ProxyOptions struct {
	Site    string
	Command []string
}

func TerraformProxy(ctx context.Context, cfg *config.MachConfig, locations map[string]string, options *ProxyOptions) error {
	for i := range cfg.Sites {
		site := cfg.Sites[i]

		if options.Site != "" && site.Identifier != options.Site {
			continue
		}

		err := RunTerraform(ctx, locations[site.Identifier], options.Command...)
		if err != nil {
			return err
		}
	}
	return nil
}

func TerraformApplySite(ctx context.Context, cfg *config.MachConfig, site *config.SiteConfig, path string, options *ApplyOptions) error {
	if !options.Reuse {
		if err := RunTerraform(ctx, path, "init"); err != nil {
			return err
		}
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
	planFilename, err := TerraformPlanDetect(path)
	if err != nil {
		return err
	}
	if planFilename != "" {
		cmd = append(cmd, planFilename)
	}

	return RunTerraform(ctx, path, cmd...)
}

func TerraformPlanDetect(path string) (string, error) {
	filename, err := GeneratePlanName(path)
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(path, filename)
	if _, err := os.Stat(filePath); err == nil {
		return filename, nil
	}
	return "", nil
}
