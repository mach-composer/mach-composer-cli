package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return version information of the mach-composer cli",
	Run: func(cmd *cobra.Command, args []string) {
		if len(commit) >= 7 {
			commit = commit[:7]
		}
		fmt.Printf("mach-composer %s (%s) - %s\n", version, commit, date)
	},
}
