package cmd

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/cmd/cloudcmd"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
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
				cli.PrintExitError(err.Error())
			}

			output, err := cmd.Flags().GetString("output")
			if err != nil {
				cli.PrintExitError(err.Error())
			}

			stripLogs, err := cmd.Flags().GetBool("strip-logs")
			if err != nil {
				cli.PrintExitError(err.Error())
			}

			var w io.Writer
			switch output {
			case string(cli.OutputTypeJSON):
				w = os.Stdout
			case string(cli.OutputTypeConsole):
				var partsExclude []string
				if stripLogs {
					partsExclude = []string{zerolog.LevelFieldName, utils.IdentifierFieldName}
				}

				var fieldsExclude []string
				if !verbose {
					fieldsExclude = []string{
						utils.ArgsFieldName, utils.CommandFieldName, utils.CwdFieldName, utils.IdentifierFieldName,
						utils.NameName, utils.ModuleName, utils.TimestampMain,
					}
				}

				w = zerolog.ConsoleWriter{
					Out:           os.Stdout,
					PartsOrder:    []string{zerolog.LevelFieldName, utils.IdentifierFieldName, zerolog.MessageFieldName},
					FieldsExclude: fieldsExclude,
					PartsExclude:  partsExclude,
				}
			default:
				cli.PrintExitError("unknown output type: %s", output)
			}
			ctx = cli.ContextWithOutput(ctx, cli.OutputType(output))
			ctx = cli.ContextWithLogWriter(ctx, w)

			var logger = zerolog.New(w).With().Str(utils.IdentifierFieldName, utils.FormatIdentifier("main")).Timestamp().Logger()

			if verbose {
				logger = logger.Level(zerolog.TraceLevel)
			} else if quiet {
				logger = logger.Level(zerolog.ErrorLevel)
			} else {
				logger = logger.Level(zerolog.InfoLevel)
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
	RootCmd.PersistentFlags().Bool("strip-logs", false, "Strip all context from the logs")
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
