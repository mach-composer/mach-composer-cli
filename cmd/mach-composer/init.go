package main

import (
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/generator"
	"github.com/mach-composer/mach-composer-cli/internal/runner"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize site directories Terraform files.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessCommonFlags()
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
}

func initFunc(cmd *cobra.Command, _ []string) error {
	cfg := loadConfig(cmd, true)
	defer cfg.Close()
	ctx := cmd.Context()

	dg, err := dependency.ToDeploymentGraph(cfg, commonFlags.outputPath)
	if err != nil {
		return err
	}

	err = generator.Write(ctx, cfg, dg, nil)
	if err != nil {
		return err
	}

	return runner.TerraformInit(ctx, cfg, dg, nil)
}
