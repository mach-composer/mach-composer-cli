package main

import (
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

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

	RunE: func(cmd *cobra.Command, args []string) error {
		return showPlanFunc(cmd, args)
	},
}

func init() {
	registerGenerateFlags(showPlanCmd)
	showPlanCmd.Flags().BoolVarP(&showPlanFlags.noColor, "no-color", "", false, "Disable color output")
}

func showPlanFunc(cmd *cobra.Command, _ []string) error {
	cfg := loadConfig(cmd, true)
	defer cfg.Close()
	ctx := cmd.Context()

	generateFlags.ValidateSite(cfg)

	dg, err := dependency.ToDeploymentGraph(cfg)
	if err != nil {
		return err
	}

	if generateFlags.siteName != "" {
		log.Warn().Msgf("Site option not implemented")
	}

	return runner.TerraformShow(ctx, cfg, dg, &runner.ShowPlanOptions{
		NoColor: showPlanFlags.noColor,
	})
}
