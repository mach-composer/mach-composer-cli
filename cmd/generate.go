package cmd

import "github.com/spf13/cobra"

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the Terraform files.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// if err := someFunc(); err != nil {
		// 	return err
		// }
		return nil
	},
}
