package cli

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
)

// TerraformHook is a hook that will parse the log output from terraform if set
// to json. Otherwise, the existing output is returned
type TerraformHook struct {
	identifier string
}

func NewTerraformHook(identifier string) *TerraformHook {
	return &TerraformHook{
		identifier: identifier,
	}
}

type LogLine struct {
	Level     string                 `mapstructure:"@level"`
	Message   string                 `mapstructure:"@message"`
	Module    string                 `mapstructure:"@module"`
	Timestamp string                 `mapstructure:"@timestamp"`
	Remainder map[string]interface{} `mapstructure:",remain"`
}

func (g *TerraformHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	if OutputFromContext(e.GetCtx()) == OutputTypeJSON {
		logger := log.Logger.With().Str("identifier", g.identifier).Logger()

		for _, line := range strings.Split(message, "\n") {
			if line == "" {
				continue
			}

			var logLine = LogLine{
				Remainder: make(map[string]interface{}),
			}

			data := make(map[string]interface{})
			err := json.Unmarshal([]byte(line), &data)
			if err != nil {
				logger.Err(err).Str("line", line).Msg("Failed to unmarshal log line")
				continue
			}

			err = mapstructure.Decode(data, &logLine)
			if err != nil {
				logger.Err(err).Str("line", line).Any("data", data).Msg("Failed to decode log line")
				continue
			}

			level, err = zerolog.ParseLevel(logLine.Level)
			if err != nil {
				logger.Err(err).Str("line", line).Str("logLevel", logLine.Level).Msg("Failed to parse log level")
				continue
			}

			logger.WithLevel(level).Fields(logLine.Remainder).Msg(logLine.Message)
		}

		e.Discard()
		return
	}

	log.WithLevel(level).Msg(message)
	e.Discard()
	return
}
