package cli

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"strings"
)

type LogLine struct {
	Level     string                 `mapstructure:"@level"`
	Message   string                 `mapstructure:"@message"`
	Module    string                 `mapstructure:"@module"`
	Timestamp string                 `mapstructure:"@timestamp"`
	Remainder map[string]interface{} `mapstructure:",remain"`
}

func ParseTerraformJsonOutput(message string) ([]LogLine, error) {
	var logLines []LogLine

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
			return nil, err
		}

		err = mapstructure.Decode(data, &logLine)
		if err != nil {
			return nil, err
		}

		logLines = append(logLines, logLine)
	}

	return logLines, nil
}
