package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	rootCmd = &cobra.Command{
		Use:   "mach-composer",
		Short: "MACH composer is an orchestration tool for modern MACH ecosystems",
		Long: `MACH composer is a framework that you use to orchestrate and` +
			`extend modern digital commerce & experience platforms, based on MACH ` +
			`technologies and cloud native services.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				logrus.SetLevel(logrus.DebugLevel)
			}
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "", false, "Verbose output.")

	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(planCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(versionCmd)
}
