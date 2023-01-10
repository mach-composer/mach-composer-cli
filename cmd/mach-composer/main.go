package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/labd/mach-composer/cmd/mach-composer/cloudcmd"
	"github.com/labd/mach-composer/internal/cli"
)

var (
	name    = "mach-composer"
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	rootCmd = &cobra.Command{
		Use:   name,
		Short: "MACH composer is an orchestration tool for modern MACH ecosystems",
		Long: `MACH composer is a framework that you use to orchestrate and ` +
			`extend modern digital commerce & experience platforms, based on MACH ` +
			`technologies and cloud native services.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			verbose, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				panic(err)
			}

			logger := zerolog.New(os.Stdout)
			if verbose {
				logger = logger.Level(zerolog.TraceLevel)
			} else {
				logger = logger.Level(zerolog.InfoLevel)
			}
			logger = logger.Output(cli.NewConsoleWriter())
			log.Logger = logger
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "", false, "Verbose output.")
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(cloudcmd.CloudCmd)
	rootCmd.AddCommand(componentsCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(planCmd)
	rootCmd.AddCommand(schemaCmd)
	rootCmd.AddCommand(showPlanCmd)
	rootCmd.AddCommand(sitesCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(terraformCmd)
	rootCmd.AddCommand(versionCmd)
}
