package main

import (
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/generator"
	"github.com/mach-composer/mach-composer-cli/internal/runner"
)

var planFlags struct {
	reuse      bool
	components []string
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
	planCmd.Flags().BoolVarP(&planFlags.reuse, "reuse", "", false, "Suppress a terraform init for improved speed (not recommended for production usage)")
	planCmd.Flags().StringArrayVarP(&planFlags.components, "component", "c", []string{}, "")
}

func planFunc(cmd *cobra.Command, args []string) error {
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

	return runner.TerraformPlan(ctx, cfg, dg, &runner.PlanOptions{
		Reuse: planFlags.reuse,
		Site:  generateFlags.siteName,
	})
}
