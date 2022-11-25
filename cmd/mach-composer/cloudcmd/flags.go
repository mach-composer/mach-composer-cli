package cloudcmd

import (
	"github.com/spf13/cobra"
)

func registerContextFlags(cmd *cobra.Command) {
	cmd.Flags().String("organization", "", "Organization key")
	cmd.Flags().String("project", "", "Project key")
	Must(cmd.MarkFlagRequired("organization"))
	Must(cmd.MarkFlagRequired("project"))
}
