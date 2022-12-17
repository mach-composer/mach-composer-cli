package cli

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)


	}


func PrintExitError(summary string, detail ...string) {
	log.Error().
		Str("details", strings.Join(detail, "\n")).
		Msg(summary)
	os.Exit(1)
}
