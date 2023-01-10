package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/labd/mach-composer/internal/config"
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

		err := TerraformPlanSite(ctx, cfg, &site, locations[site.Identifier], options)
		if err != nil {
			return err
		}
	}
	return nil
}

func TerraformPlanSite(ctx context.Context, cfg *config.MachConfig, site *config.SiteConfig, path string, options *PlanOptions) error {
	log.Debug().Msgf("Running terraform plan for site %s", site.Identifier)

	if !options.Reuse {
		if err := RunTerraform(ctx, path, "init"); err != nil {
			return err
		}
	}
	cmd := []string{"plan"}

	for _, component := range options.Components {
		cmd = append(cmd, fmt.Sprintf("-target=module.%s", component))
	}

	filename, err := generatePlanName(path)
	if err != nil {
		return err
	}
	cmd = append(cmd, fmt.Sprintf("-out=%s", filename))

	return RunTerraform(ctx, path, cmd...)
}

func generatePlanName(path string) (string, error) {
	siteHash, err := getFileHash(path)
	if err != nil {
		return "", err
	}
	result := fmt.Sprintf("%s.tfplan", siteHash[:7])
	return result, nil
}

func findExistingPlan(path string) (string, error) {
	filename, err := generatePlanName(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("No terraform files found. Did you run `mach-composer generate`?")
		}
		return "", err
	}
	filePath := filepath.Join(path, filename)
	if _, err := os.Stat(filePath); err == nil {
		return filename, nil
	}
	return "", nil
}
