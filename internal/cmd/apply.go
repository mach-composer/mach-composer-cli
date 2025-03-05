package cmd

import (
	"github.com/mach-composer/mach-composer-cli/internal/batcher"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/mach-composer/mach-composer-cli/internal/hash"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/generator"
	"github.com/mach-composer/mach-composer-cli/internal/runner"
)

var applyFlags struct {
	forceInit             bool
	autoApprove           bool
	destroy               bool
	github                bool
	bufferLogs            bool
	components            []string
	numWorkers            int
	ignoreChangeDetection bool
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the configuration.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessCommonFlags(cmd)
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		return applyFunc(cmd, args)
	},
}

func init() {
	registerCommonFlags(applyCmd)
	applyCmd.Flags().BoolVarP(&applyFlags.forceInit, "force-init", "", false, "Force terraform initialization. By default mach-composer will reuse existing terraform resources")
	applyCmd.Flags().BoolVarP(&applyFlags.autoApprove, "auto-approve", "", false, "Suppress a terraform init for improved speed (not recommended for production usage)")
	applyCmd.Flags().BoolVarP(&applyFlags.destroy, "destroy", "", false, "Destroy option is a convenient way to destroy all remote objects managed by this mach config")
	applyCmd.Flags().StringArrayVarP(&applyFlags.components, "component", "c", nil, "")
	applyCmd.Flags().BoolVarP(&applyFlags.ignoreChangeDetection, "ignore-change-detection", "", false, "Ignore change detection to run even if the components are considered up to date")
	applyCmd.Flags().BoolVarP(&applyFlags.github, "github", "g", false, "Whether logs should be decorated with github-specific formatting")
	applyCmd.Flags().BoolVarP(&applyFlags.bufferLogs, "buffer", "b", false, "Whether logs should be buffered and printed at the end of the run")
}

func applyFunc(cmd *cobra.Command, _ []string) error {
	if len(applyFlags.components) > 0 {
		log.Warn().Msgf("Components option not implemented")
	}

	if applyFlags.github && !applyFlags.bufferLogs {
		log.Warn().Msg("Github flag is only supported with buffer flag")
	}

	cfg := loadConfig(cmd, true)
	defer cfg.Close()
	ctx := cmd.Context()

	dg, err := graph.ToDeploymentGraph(cfg, commonFlags.outputPath)
	if err != nil {
		return err
	}

	// Note that we do this in multiple passes to minimize ending up with
	// half broken runs. We could in the future also run some parts in parallel

	err = generator.Write(ctx, cfg, dg, nil)
	if err != nil {
		return err
	}

	r := runner.NewGraphRunner(
		batcher.NaiveBatchFunc(),
		hash.Factory(cfg),
		commonFlags.workers,
	)

	return r.TerraformApply(ctx, dg, &runner.ApplyOptions{
		ForceInit:             applyFlags.forceInit,
		Destroy:               applyFlags.destroy,
		AutoApprove:           applyFlags.autoApprove,
		IgnoreChangeDetection: applyFlags.ignoreChangeDetection,
		BufferLogs:            applyFlags.bufferLogs,
		Github:                applyFlags.github,
	})
}
