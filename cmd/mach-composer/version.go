package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return version information of the mach-composer cli",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(commit) >= 7 {
			commit = commit[:7]
		}
		fmt.Printf("mach-composer %s (%s) - %s\n", version, commit, date)
		return nil
	},
}
