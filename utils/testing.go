package utils

import (
	"strings"

	"github.com/lithammer/dedent"
)

func TrimIndent(data string) string {
	return dedent.Dedent(strings.ReplaceAll(data, "\t", "    "))
}
