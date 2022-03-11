package cmd

import (
	"github.com/labd/mach-composer-go/config"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the Terraform files.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := generateFunc(); err != nil {
			return err
		}
		return nil
	},
}

func generateFunc() error {
	data := []byte(``)
	config.Parse(data)

	return nil

}
