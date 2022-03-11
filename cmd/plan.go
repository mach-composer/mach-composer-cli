package cmd

import "github.com/spf13/cobra"

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Output the deploy plan.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// if err := someFunc(); err != nil {
		// 	return err
		// }
		return nil
	},
}
