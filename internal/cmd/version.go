package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/cli"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return version information of the mach-composer cli",
	Run: func(cmd *cobra.Command, args []string) {
		md := cli.GetVersionMetadata()
		fmt.Println(md.String())
	},
}
