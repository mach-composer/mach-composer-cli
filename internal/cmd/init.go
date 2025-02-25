package cmd

import (
	"github.com/mach-composer/mach-composer-cli/internal/batcher"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/mach-composer/mach-composer-cli/internal/hash"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/generator"
	"github.com/mach-composer/mach-composer-cli/internal/runner"
)

var initFlags struct {
	github     bool
	bufferLogs bool
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize site directories Terraform files.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessCommonFlags(cmd)
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		cli.DeprecationWarning(&cli.DeprecationOptions{
			Message: "the init command will change in the next version. For initializing terraform please use 'mach-composer terraform init'.",
		})
		return initFunc(cmd, args)
	},
}

func init() {
	registerCommonFlags(initCmd)
	initCmd.Flags().BoolVarP(&initFlags.github, "github", "g", false, "Whether logs should be decorated with github-specific formatting")
	initCmd.Flags().BoolVarP(&initFlags.bufferLogs, "buffer", "b", false, "Whether logs should be buffered and printed at the end of the run")
}

func initFunc(cmd *cobra.Command, _ []string) error {
	if initFlags.github && !initFlags.bufferLogs {
		log.Warn().Msg("Github flag is only supported with buffer flag")
	}

	cfg := loadConfig(cmd, true)
	defer cfg.Close()
	ctx := cmd.Context()

	dg, err := graph.ToDeploymentGraph(cfg, commonFlags.outputPath)
	if err != nil {
		return err
	}

	err = generator.Write(ctx, cfg, dg, nil)
	if err != nil {
		return err
	}

	r := runner.NewGraphRunner(
		batcher.NaiveBatchFunc(),
		hash.Factory(cfg),
		commonFlags.workers,
	)

	return r.TerraformInit(ctx, dg, &runner.InitOptions{
		BufferLogs: initFlags.bufferLogs,
		Github:     initFlags.github,
	})
}
