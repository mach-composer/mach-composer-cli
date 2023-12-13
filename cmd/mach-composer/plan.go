package main

import (
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/generator"
	"github.com/mach-composer/mach-composer-cli/internal/runner"
)

var planFlags struct {
	reuse      bool
	components []string
	lock       bool
}

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Plan the configuration.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		return planFunc(cmd, args)
	},
}

func init() {
	registerGenerateFlags(planCmd)
	planCmd.Flags().BoolVarP(&planFlags.reuse, "reuse", "", false,
		"Suppress a terraform init for improved speed (not recommended for production usage)")
	planCmd.Flags().StringArrayVarP(&planFlags.components, "component", "c", []string{}, "")
	planCmd.Flags().BoolVarP(&planFlags.lock, "lock", "", true,
		"Acquire a lock on the state file before running terraform plan")
}

func planFunc(cmd *cobra.Command, _ []string) error {
	cfg := loadConfig(cmd, true)
	defer cfg.Close()
	ctx := cmd.Context()

	generateFlags.ValidateSite(cfg)

	dg, err := dependency.ToDeploymentGraph(cfg)
	if err != nil {
		return err
	}

	err = generator.Write(ctx, cfg, dg, &generator.GenerateOptions{
		OutputPath: generateFlags.outputPath,
		Site:       generateFlags.siteName,
	})
	if err != nil {
		return err
	}

	if len(planFlags.components) > 0 {
		log.Warn().Msgf("Components option not implemented")
	}
	if generateFlags.siteName != "" {
		log.Warn().Msgf("Site option not implemented")
	}

	return runner.TerraformPlan(ctx, cfg, dg, &runner.PlanOptions{
		Reuse: planFlags.reuse,
		Lock:  planFlags.lock,
	})
}
