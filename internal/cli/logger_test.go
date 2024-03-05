package cli

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
)

type LogMessage struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

type LogSink struct {
	logs []LogMessage
}

func (l *LogSink) Write(p []byte) (n int, err error) {
	var msg = LogMessage{}
	_ = json.Unmarshal(p, &msg)

	l.logs = append(l.logs, msg)
	return len(p), nil
}

func (l *LogSink) Index(i int) LogMessage {
	return l.logs[i]
}

func SetSinkLogger(sink *LogSink) func() {
	oldLogger := log.Logger
	log.Logger = log.Output(sink)

	return func() {
		log.Logger = oldLogger
	}
}
