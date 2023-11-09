package main

import (
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/generator"
	"github.com/mach-composer/mach-composer-cli/internal/runner"
)

var applyFlags struct {
	reuse       bool
	autoApprove bool
	destroy     bool
	components  []string
	numWorkers  int
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the configuration.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		return applyFunc(cmd, args)
	},
}

func init() {
	registerGenerateFlags(applyCmd)
	applyCmd.Flags().BoolVarP(&applyFlags.reuse, "reuse", "", false, "Suppress a terraform init for improved speed (not recommended for production usage)")
	applyCmd.Flags().BoolVarP(&applyFlags.autoApprove, "auto-approve", "", false, "Suppress a terraform init for improved speed (not recommended for production usage)")
	applyCmd.Flags().BoolVarP(&applyFlags.destroy, "destroy", "", false, "Destroy option is a convenient way to destroy all remote objects managed by this mach config")
	applyCmd.Flags().StringArrayVarP(&applyFlags.components, "component", "c", []string{}, "")
	applyCmd.Flags().IntVarP(&applyFlags.numWorkers, "workers", "w", 1, "The number of workers to use")
}

func applyFunc(cmd *cobra.Command, args []string) error {
	cfg := loadConfig(cmd, true)
	defer cfg.Close()
	ctx := cmd.Context()

	generateFlags.ValidateSite(cfg)

	dg, err := dependency.ToDeploymentGraph(cfg)
	if err != nil {
		return err
	}

	// Note that we do this in multiple passes to minimize ending up with
	// half broken runs. We could in the future also run some parts in parallel

	err = generator.Write(ctx, cfg, dg, &generator.GenerateOptions{
		OutputPath: generateFlags.outputPath,
		Site:       generateFlags.siteName,
	})
	if err != nil {
		return err
	}

	return runner.TerraformApply(ctx, cfg, dg, &runner.ApplyOptions{
		Destroy:     applyFlags.destroy,
		Reuse:       applyFlags.reuse,
		AutoApprove: applyFlags.autoApprove,
		Site:        generateFlags.siteName,
		Components:  applyFlags.components,
	})
}
