package cmd

import "github.com/spf13/cobra"

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the configuration.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// if err := someFunc(); err != nil {
		// 	return err
		// }
		return nil
	},
}
