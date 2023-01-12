package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/labd/mach-composer/internal/generator"
	"github.com/labd/mach-composer/internal/runner"
)

var showPlanFlags struct {
	noColor bool
}

var showPlanCmd = &cobra.Command{
	Use:   "show-plan",
	Short: "Show the planned configuration.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},

	Run: func(cmd *cobra.Command, args []string) {
		handleError(showPlanFunc(cmd.Context(), args))
	},
}

func init() {
	registerGenerateFlags(showPlanCmd)
	showPlanCmd.Flags().BoolVarP(&showPlanFlags.noColor, "no-color", "", false, "Disable color output")
}

func showPlanFunc(ctx context.Context, args []string) error {
	cfg := loadConfig(ctx, true)
	defer cfg.Close()

	generateFlags.ValidateSite(cfg)

	paths := generator.FileLocations(cfg, &generator.GenerateOptions{
		OutputPath: generateFlags.outputPath,
		Site:       generateFlags.siteName,
	})

	return runner.TerraformShow(ctx, cfg, paths, &runner.ShowPlanOptions{
		NoColor: showPlanFlags.noColor,
		Site:    generateFlags.siteName,
	})
}
