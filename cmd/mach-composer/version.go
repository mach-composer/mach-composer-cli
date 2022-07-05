package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return version information of the mach-composer cli",
	RunE: func(cmd *cobra.Command, args []string) error {
		if commit != "" && date != "" {
			fmt.Printf("v%s (%s) - %s\n", version, commit[:7], date)
		} else {
			if version != "development" {
				version = fmt.Sprintf("v%s", version)
			}
			fmt.Println(version)
		}
		return nil
	},
}
