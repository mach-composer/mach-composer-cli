package cmd

import (
	"github.com/mach-composer/mach-composer-cli/internal/batcher"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/mach-composer/mach-composer-cli/internal/hash"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/runner"
)

var showPlanFlags struct {
	forceInit             bool
	noColor               bool
	ignoreChangeDetection bool
	github                bool
	bufferLogs            bool
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
	showPlanCmd.Flags().BoolVarP(&showPlanFlags.github, "github", "g", false, "Whether logs should be decorated with github-specific formatting")
	showPlanCmd.Flags().BoolVarP(&showPlanFlags.bufferLogs, "buffer", "b", false, "Whether logs should be buffered and printed at the end of the run")
}

func showPlanFunc(cmd *cobra.Command, _ []string) error {
	if showPlanFlags.github && !showPlanFlags.bufferLogs {
		log.Warn().Msg("Github flag is only supported with buffer flag")
	}

	cfg := loadConfig(cmd, true)
	defer cfg.Close()
	ctx := cmd.Context()

	dg, err := graph.ToDeploymentGraph(cfg, commonFlags.outputPath)
	if err != nil {
		return err
	}

	b, err := batcher.Factory(cfg)
	if err != nil {
		return err
	}

	h := hash.Factory(cfg)

	r := runner.NewGraphRunner(b, h, commonFlags.workers)

	return r.TerraformShow(ctx, dg, &runner.ShowPlanOptions{
		ForceInit:             showPlanFlags.forceInit,
		NoColor:               showPlanFlags.noColor,
		IgnoreChangeDetection: showPlanFlags.ignoreChangeDetection,
		Github:                showPlanFlags.github,
		BufferLogs:            showPlanFlags.bufferLogs,
	})
}
