package cmd

import (
	"github.com/mach-composer/mach-composer-cli/internal/batcher"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/mach-composer/mach-composer-cli/internal/hash"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/runner"
)

var terraformFlags struct {
	reuse                 bool
	ignoreChangeDetection bool
}

var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Execute terraform commands directly",
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
}

func terraformFunc(cmd *cobra.Command, args []string) error {
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
	})
}
