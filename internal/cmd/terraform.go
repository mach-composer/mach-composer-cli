package cmd

import (
	"github.com/mach-composer/mach-composer-cli/internal/batcher"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/mach-composer/mach-composer-cli/internal/hash"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/runner"
)

var terraformFlags struct {
	reuse                 bool
	ignoreChangeDetection bool
	github                bool
	bufferLogs            bool
	filters               []string
}

var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Execute terraform commands directly. See [the documentation](/howto/cli/filtering-commands) for filtering options.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessCommonFlags(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return terraformFunc(cmd, args)
	},
}

func init() {
	registerCommonFlags(terraformCmd)
	terraformCmd.Flags().BoolVarP(&terraformFlags.ignoreChangeDetection, "ignore-change-detection", "", true,
		"Ignore change detection to run even if the components are considered up to date. Per default the proxy will ignore change detection")
	terraformCmd.Flags().BoolVarP(&terraformFlags.github, "github", "g", false, "Whether logs should be decorated with github-specific formatting")
	terraformCmd.Flags().BoolVarP(&terraformFlags.bufferLogs, "buffer", "b", false, "Whether logs should be buffered and printed at the end of the run")
	terraformCmd.Flags().StringArrayVarP(&terraformFlags.filters, "filter", "", nil, "Run only nodes matching the filter expression")
}

func terraformFunc(cmd *cobra.Command, args []string) error {
	if terraformFlags.github && !terraformFlags.bufferLogs {
		log.Warn().Msg("Github flag is only supported with buffer flag")
	}
	cfg := loadConfig(cmd, true)
	ctx := cmd.Context()
	defer cfg.Close()

	dg, err := graph.ToDeploymentGraph(cfg, commonFlags.outputPath)
	if err != nil {
		return err
	}

	r := runner.NewGraphRunner(
		batcher.NaiveBatchFunc(),
		hash.Factory(cfg),
		commonFlags.workers,
	)

	return r.TerraformProxy(ctx, dg, &runner.ProxyOptions{
		Command:               args,
		IgnoreChangeDetection: terraformFlags.ignoreChangeDetection,
		Github:                terraformFlags.github,
		BufferLogs:            terraformFlags.bufferLogs,
		Filters:               terraformFlags.filters,
	})
}
