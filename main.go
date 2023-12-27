package main

import (
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/cmd"
	"os"
)

//go:generate go run tools/cli-docs/main.go docs/src/reference/cli/
func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		if cmd.RootCmd.SilenceErrors {
			cli.HandleErr(err)
		}
		os.Exit(1)
	}
}
