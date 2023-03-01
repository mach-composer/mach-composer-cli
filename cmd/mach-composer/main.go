package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/labd/mach-composer/cmd/mach-composer/cloudcmd"
	"github.com/labd/mach-composer/internal/cli"
)

var (
	rootCmd = &cobra.Command{
		Use:   "mach-composer",
		Short: "MACH composer is an orchestration tool for modern MACH ecosystems",
		Long: `MACH composer is a framework that you use to orchestrate and ` +
			`extend modern digital commerce & experience platforms, based on MACH ` +
			`technologies and cloud native services.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.Root().SilenceUsage = true
			cmd.Root().SilenceErrors = true

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

			ctx := logger.WithContext(cmd.Context())
			ctx, cancel := context.WithCancel(ctx)

			// Register a signal handler to cancel the current context
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)

			go func() {
				select {
				case <-c:
					log.Info().Msg("Exiting...")
					cancel()
				case <-ctx.Done():
				}
			}()

			cmd.SetContext(ctx)
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		if rootCmd.SilenceErrors {
			cli.HandleErr(err)
		}
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
