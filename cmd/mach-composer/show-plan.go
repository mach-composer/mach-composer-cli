package main

import (
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/generator"
	"github.com/mach-composer/mach-composer-cli/internal/runner"
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
		handleError(showPlanFunc(cmd, args))
	},
}

func init() {
	registerGenerateFlags(showPlanCmd)
	showPlanCmd.Flags().BoolVarP(&showPlanFlags.noColor, "no-color", "", false, "Disable color output")
}

func showPlanFunc(cmd *cobra.Command, args []string) error {
	cfg := loadConfig(cmd, true)
	defer cfg.Close()
	ctx := cmd.Context()

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
