package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/labd/mach-composer/internal/generator"
	"github.com/labd/mach-composer/internal/runner"
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
		return handleError(planFunc(cmd.Context(), args))
	},
}

func init() {
	registerGenerateFlags(planCmd)
	planCmd.Flags().BoolVarP(&planFlags.reuse, "reuse", "", false, "Suppress a terraform init for improved speed (not recommended for production usage)")
	planCmd.Flags().StringArrayVarP(&planFlags.components, "component", "c", []string{}, "")
}

func planFunc(ctx context.Context, args []string) error {
	cfg := loadConfig(ctx, true)
	generateFlags.ValidateSite(cfg)

	paths, err := generator.WriteFiles(cfg, &generator.GenerateOptions{
		OutputPath: generateFlags.outputPath,
		Site:       generateFlags.siteName,
	})
	if err != nil {
		return err
	}

	return runner.TerraformPlan(ctx, cfg, paths, &runner.PlanOptions{
		Reuse: planFlags.reuse,
		Site:  generateFlags.siteName,
	})
}
