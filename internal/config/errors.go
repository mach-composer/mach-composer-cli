package config

import (
	"fmt"
	"strings"
)

type SyntaxError struct {
	message  string
	line     int
	column   int
	filename string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("%s on at %s at line %d:%d", e.message, e.filename, e.line, e.column)
}

type ValidationError struct {
	errors []string
}

func (e *ValidationError) Error() string {
	lines := []string{}
	for _, err := range e.errors {
		lines = append(lines, fmt.Sprintf(" - %s", err))
	}
	return fmt.Sprintf(
		"The configuration is not valid:\n%s",
		strings.Join(lines, ""))
}
