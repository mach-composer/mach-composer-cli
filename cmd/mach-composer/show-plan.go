package main

import (
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/runner"
)

var showPlanFlags struct {
	reuse   bool
	noColor bool
}

var showPlanCmd = &cobra.Command{
	Use:   "show-plan",
	Short: "Show the planned configuration.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessCommonFlags()
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
}

func showPlanFunc(cmd *cobra.Command, _ []string) error {
	cfg := loadConfig(cmd, true)
	defer cfg.Close()
	ctx := cmd.Context()

	dg, err := dependency.ToDeploymentGraph(cfg, commonFlags.outputPath)
	if err != nil {
		return err
	}

	if showPlanFlags.reuse == false {
		if err = runner.TerraformInit(ctx, cfg, dg, nil); err != nil {
			return err
		}
	} else {
		log.Info().Msgf("Reusing existing terraform state")
	}

	return runner.TerraformShow(ctx, dg, &runner.ShowPlanOptions{
		NoColor: showPlanFlags.noColor,
		Workers: commonFlags.workers,
	})
}
