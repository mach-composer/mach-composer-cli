package main

import (
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/runner"
)

var terraformFlags struct {
	reuse bool
}

var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Execute terraform commands directly",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessCommonFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return terraformFunc(cmd, args)
	},
}

func init() {
	registerCommonFlags(terraformCmd)
	terraformCmd.Flags().BoolVarP(&terraformFlags.reuse, "reuse", "", false,
		"Suppress a terraform init for improved speed (not recommended for production usage)")
}

func terraformFunc(cmd *cobra.Command, args []string) error {
	cfg := loadConfig(cmd, true)
	ctx := cmd.Context()
	defer cfg.Close()

	dg, err := dependency.ToDeploymentGraph(cfg, commonFlags.outputPath)
	if err != nil {
		return err
	}

	if terraformFlags.reuse == false {
		if err = runner.TerraformInit(ctx, cfg, dg, nil); err != nil {
			return err
		}
	} else {
		log.Info().Msgf("Reusing existing terraform state")
	}

	return runner.TerraformProxy(cmd.Context(), dg, &runner.ProxyOptions{
		Command: args,
		Workers: commonFlags.workers,
	})
}
