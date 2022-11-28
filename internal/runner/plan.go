package runner

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

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
	logrus.Debugf("Running terraform plan for site %s", site.Identifier)

	if !options.Reuse {
		if err := RunTerraform(ctx, path, "init"); err != nil {
			return err
		}
	}
	cmd := []string{"plan"}

	for _, component := range options.Components {
		cmd = append(cmd, fmt.Sprintf("-target=module.%s", component))
	}

	filename, err := GeneratePlanName(path)
	if err != nil {
		return err
	}
	cmd = append(cmd, fmt.Sprintf("-out=%s", filename))

	return RunTerraform(ctx, path, cmd...)
}

func GeneratePlanName(path string) (string, error) {
	siteHash, err := GetHash(path)
	if err != nil {
		return "", err
	}
	result := fmt.Sprintf("%s.tfplan", siteHash[:7])
	return result, nil
}
