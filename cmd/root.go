package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "mach",
		Short: "MACH composer is an orchestration tool for modern MACH ecosystems",
		Long: `MACH composer is a framework that you use to orchestrate and` +
			`extend modern digital commerce & experience platforms, based on MACH ` +
			`technologies and cloud native services.`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(planCmd)
	rootCmd.AddCommand(updateCmd)
}
