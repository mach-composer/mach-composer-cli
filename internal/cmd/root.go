package cmd

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/cmd/cloudcmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/signal"
)

var (
	RootCmd = &cobra.Command{
		Use:   "mach-composer",
		Short: "MACH composer is an orchestration tool for modern MACH ecosystems",
		Long: `MACH composer is a framework that you use to orchestrate and ` +
			`extend modern digital commerce & experience platforms, based on MACH ` +
			`technologies and cloud native services.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()

			cmd.Root().SilenceUsage = true
			cmd.Root().SilenceErrors = true

			verbose, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				cli.PrintExitError(err.Error())
			}

			quiet, err := cmd.Flags().GetBool("quiet")
			if err != nil {
				panic(err)
			}

			//Configure logger
			output, err := cmd.Flags().GetString("output")
			var w io.Writer
			switch output {
			case string(cli.OutputTypeJSON):
				w = os.Stdout
			case string(cli.OutputTypeConsole):
				w = cli.NewConsoleWriter()
			default:
				cli.PrintExitError("unknown output type: %s", output)
			}
			ctx = cli.ContextWithOutput(ctx, cli.OutputType(output))
			ctx = cli.ContextWithLogWriter(ctx, w)

			var logger = zerolog.New(w).With().Timestamp().Logger()

			if verbose {
				logger = logger.Level(zerolog.TraceLevel)
			} else if quiet {
				logger = logger.Level(zerolog.ErrorLevel)
			} else {
				logger = logger.Level(zerolog.InfoLevel)
			}

			github, err := cmd.Flags().GetBool("github")
			if err != nil {
				panic(err)
			}
			if github {
				ctx = cli.ContextWithGithubCI(ctx)
			}

			//Load logger into context and global logger
			ctx = logger.WithContext(ctx)
			log.Logger = logger
			// Register a signal handler to cancel the current context
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)

			ctx, cancel := context.WithCancel(ctx)

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

func init() {
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output. This is equal to setting log levels to debug and higher")
	RootCmd.PersistentFlags().BoolP("quiet", "q", false, "Quiet output. This is equal to setting log levels to error and higher")
	RootCmd.PersistentFlags().String("output", "console", "The output type. One of: console, json")
	RootCmd.PersistentFlags().BoolP("github", "g", false, "Whether logs should be decorated with github-specific formatting")
	RootCmd.AddCommand(applyCmd)
	RootCmd.AddCommand(cloudcmd.CloudCmd)
	RootCmd.AddCommand(componentsCmd)
	RootCmd.AddCommand(generateCmd)
	RootCmd.AddCommand(initCmd)
	RootCmd.AddCommand(planCmd)
	RootCmd.AddCommand(schemaCmd)
	RootCmd.AddCommand(showPlanCmd)
	RootCmd.AddCommand(sitesCmd)
	RootCmd.AddCommand(updateCmd)
	RootCmd.AddCommand(terraformCmd)
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(graphCmd)
	RootCmd.AddCommand(validateCmd)
}
