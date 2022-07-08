package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return version information of the mach-composer cli",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("mach-composer %s (%s) - %s\n", version, commit[:7], date)
		return nil
	},
}
