package cli

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/labd/mach-composer/internal/utils"
)

type DeprecationOptions struct {
	Site      string
	Component string
	Message   string
	Details   string
	Version   string
}

func DeprecationWarning(opts *DeprecationOptions) {
	if opts.Details != "" {
		opts.Details = utils.TrimIndent(opts.Details)
	}

	log.Warn().
		Str("site", opts.Site).
		Str("component", opts.Component).
		Str("deprecation", opts.Version).
		Str("details", opts.Details).
		Msg(opts.Message)
}

func PrintExitError(summary string, detail ...string) {
	log.Error().
		Str("details", strings.Join(detail, "\n")).
		Msg(summary)
	os.Exit(1)
}
