package cmd

import (
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/runner"
)

var showPlanFlags struct {
	reuse                 bool
	noColor               bool
	ignoreChangeDetection bool
}

var showPlanCmd = &cobra.Command{
	Use:   "show-plan",
	Short: "Show the planned configuration.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessCommonFlags(cmd)
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		return showPlanFunc(cmd, args)
	},
}

func init() {
	registerCommonFlags(showPlanCmd)
	showPlanCmd.Flags().BoolVarP(&showPlanFlags.noColor, "no-color", "", false, "Disable color output")
	showPlanCmd.Flags().BoolVarP(&showPlanFlags.reuse, "reuse", "", false,
		"Suppress a terraform init for improved speed (not recommended for production usage)")
	showPlanCmd.Flags().BoolVarP(&showPlanFlags.ignoreChangeDetection, "ignore-change-detection", "", false,
		"Ignore change detection to run even if the components are considered up to date")
}

func showPlanFunc(cmd *cobra.Command, _ []string) error {
	cfg := loadConfig(cmd, true)
	defer cfg.Close()
	ctx := cmd.Context()

	dg, err := graph.ToDeploymentGraph(cfg, commonFlags.outputPath)
	if err != nil {
		return err
	}

	b := runner.NewGraphRunner(commonFlags.workers)

	if err = checkReuse(ctx, dg, b, applyFlags.reuse); err != nil {
		return err
	}

	return b.TerraformShow(ctx, dg, &runner.ShowPlanOptions{
		NoColor:               showPlanFlags.noColor,
		IgnoreChangeDetection: showPlanFlags.ignoreChangeDetection,
	})
}
