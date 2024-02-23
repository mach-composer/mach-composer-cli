package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/rs/zerolog"

	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type ConsoleWriter struct {
}

func NewConsoleWriter() ConsoleWriter {
	w := ConsoleWriter{}

	return w
}

func (w ConsoleWriter) Write(p []byte) (n int, err error) {
	event := map[string]any{}

	d := json.NewDecoder(bytes.NewReader(p))
	d.UseNumber()
	err = d.Decode(&event)
	if err != nil {
		return n, fmt.Errorf("cannot decode event: %s", err)
	}

	var message strings.Builder
	if msg, ok := event["message"].(string); ok {
		message.WriteString(msg)
	}

	var extraFields strings.Builder
	for key, value := range event {
		extraFields.WriteString(fmt.Sprintf("  %s=%s\n", key, value))
	}

	if level, ok := event["level"].(string); ok {
		switch level {
		case zerolog.LevelTraceValue:
			c := color.New(color.FgWhite, color.Faint)
			_, _ = c.Println(message.String())
			_, _ = c.Println(extraFields.String())

		case zerolog.LevelDebugValue:
			c := color.New(color.FgWhite, color.Faint)
			_, _ = c.Println(message.String())
			_, _ = c.Println(extraFields.String())

		case zerolog.LevelInfoValue:
			c := color.New(color.Bold)
			_, _ = c.Println(message.String())

		case zerolog.LevelWarnValue:
			c := color.New(color.FgYellow, color.Bold)
			_, _ = fmt.Fprintln(os.Stderr, c.Sprint("|"))
			_, _ = fmt.Fprintln(os.Stderr, c.Sprint("| Warning:"), message.String())
			if details, ok := event["details"].(string); ok {
				printDetails(os.Stderr, details, c)
			}
			for key, value := range event {
				if key == "message" || key == "level" || key == "error" || key == "details" {
					continue
				}
				_, _ = fmt.Fprintln(os.Stderr, c.Sprint("| "), fmt.Sprintf("%s=%s", key, value))
			}
			_, _ = fmt.Fprintln(os.Stderr, c.Sprint("|"))

		case zerolog.LevelErrorValue:
			c := color.New(color.FgRed, color.Bold)
			_, _ = fmt.Fprintln(os.Stderr, c.Sprint("|"))
			_, _ = fmt.Fprintln(os.Stderr, c.Sprint("| Error:"), message.String())
			if details, ok := event["details"].(string); ok {
				printDetails(os.Stderr, details, c)
			}
			for key, value := range event {
				if key == "message" || key == "level" || key == "error" || key == "details" {
					continue
				}
				_, _ = fmt.Fprintln(os.Stderr, c.Sprint("| "), fmt.Sprintf("%s=%s", key, value))
			}
			_, _ = fmt.Fprintln(os.Stderr, c.Sprint("|"))
		}
	}

	return len(p), nil
}

func printDetails(dst io.Writer, detail string, c *color.Color) {
	if detail == "" {
		return
	}
	white := color.New(color.FgWhite, color.Bold).SprintFunc()

	line := strings.TrimSpace(utils.TrimIndent(detail))
	parts := strings.Split(line, "\n")
	_, _ = fmt.Fprintln(dst, c.Sprint("|"))
	for _, line := range parts {
		_, _ = fmt.Fprintln(dst, c.Sprint("|"), white(line))
	}
}
