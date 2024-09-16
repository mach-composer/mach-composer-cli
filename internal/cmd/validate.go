package cmd

import (
	"github.com/mach-composer/mach-composer-cli/internal/batcher"
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/mach-composer/mach-composer-cli/internal/hash"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"

	"github.com/mach-composer/mach-composer-cli/internal/generator"
	"github.com/mach-composer/mach-composer-cli/internal/runner"
)

var validateFlags struct {
	validationPath string
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the generated terraform configuration.",
	Long: "This command validates the generated terraform configuration. It will check the provided configuration file " +
		"for any errors, and will run `terraform validate` on the generated configuration. This will check for any " +
		"syntax errors in the generated configuration without accessing the actual infrastructure.\n\n" +
		"By default, the generated configuration is stored in the `validations` directory in the current " +
		"working directory. This can be changed by providing the `--validation-path` flag.\n\n" +
		"See [the terraform validation docs](https://www.terraform.io/docs/commands/validate.html) for more " +
		"information on `terraform validate`.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessCommonFlags(cmd)
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		return validateFunc(cmd, args)
	},
}

func init() {
	registerCommonFlags(validateCmd)
	validateCmd.Flags().StringVarP(&validateFlags.validationPath, "validation-path", "", "validations",
		"Directory path to store files required for configuration validation.")

	if path.IsAbs(validateFlags.validationPath) == false {
		var err error
		value, err := os.Getwd()
		if err != nil {
			cli.PrintExitError("failed to get current working directory")
		}
		validateFlags.validationPath = filepath.Join(value, validateFlags.validationPath)
	}
}

func validateFunc(cmd *cobra.Command, _ []string) error {
	cfg := loadConfig(cmd, true)
	defer cfg.Close()
	ctx := cmd.Context()

	dg, err := graph.ToDeploymentGraph(cfg, validateFlags.validationPath)
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

	return r.TerraformValidate(ctx, dg)
}
