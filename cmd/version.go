package cmd

import (
	"fmt"

	"github.com/labd/mach-composer/utils"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return version information of the mach-composer cli",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("mach-composer %s (%s) - %s\n", utils.Version, utils.Commit[:7], utils.Date)
		return nil
	},
}
