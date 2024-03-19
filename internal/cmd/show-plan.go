package cmd

import (
	"github.com/mach-composer/mach-composer-cli/internal/batcher"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/mach-composer/mach-composer-cli/internal/hash"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/runner"
)

var showPlanFlags struct {
	forceInit             bool
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
	showPlanCmd.Flags().BoolVarP(&showPlanFlags.forceInit, "force-init", "", false, "Force terraform initialization. By default mach-composer will reuse existing terraform resources")
	showPlanCmd.Flags().BoolVarP(&showPlanFlags.noColor, "no-color", "", false, "Disable color output")
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

	r := runner.NewGraphRunner(
		batcher.NaiveBatchFunc(),
		hash.Factory(cfg),
		commonFlags.workers,
	)

	return r.TerraformShow(ctx, dg, &runner.ShowPlanOptions{
		ForceInit:             showPlanFlags.forceInit,
		NoColor:               showPlanFlags.noColor,
		IgnoreChangeDetection: showPlanFlags.ignoreChangeDetection,
	})
}
