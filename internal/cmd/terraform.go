package cmd

import (
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/runner"
)

var terraformFlags struct {
	reuse bool
	force bool
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
	terraformCmd.Flags().BoolVarP(&terraformFlags.reuse, "reuse", "", false,
		"Suppress a terraform init for improved speed (not recommended for production usage)")
	terraformCmd.Flags().BoolVarP(&terraformFlags.force, "force", "", false, "Force the terraform command to run even if the components are considered up to date")
}

func terraformFunc(cmd *cobra.Command, args []string) error {
	cfg := loadConfig(cmd, true)
	ctx := cmd.Context()
	defer cfg.Close()

	dg, err := graph.ToDeploymentGraph(cfg, commonFlags.outputPath)
	if err != nil {
		return err
	}

	b := runner.NewGraphRunner(commonFlags.workers)

	if err = checkReuse(ctx, dg, b, applyFlags.reuse); err != nil {
		return err
	}

	return b.TerraformProxy(cmd.Context(), dg, &runner.ProxyOptions{
		Command: args,
		Force:   terraformFlags.force,
	})
}
