package main

import "github.com/spf13/cobra"

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize site directories Terraform files.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// if err := someFunc(); err != nil {
		// 	return err
		// }
		return nil
	},
}
