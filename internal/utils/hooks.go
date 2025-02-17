package utils

import (
	"github.com/rs/zerolog"
	"strings"
)

type StdHook struct {
	Logger zerolog.Logger
}

func (s StdHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	defer e.Discard()
	if level != zerolog.NoLevel {
		e.Send()
	}

	mParts := strings.Split(message, "\n")
	for _, m := range mParts {
		s.Logger.Info().Msg(m)
	}
}
