package main

import (
	"github.com/mach-composer/mach-composer-cli/internal/cmd"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra/doc"
	"os"
)

func main() {
	dir := os.Args[1]
	if dir == "" {
		log.Fatal().Msg("Please provide a directory to write the docs to")
	}

	cmd.RootCmd.DisableAutoGenTag = true

	err := doc.GenMarkdownTree(cmd.RootCmd, dir)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
