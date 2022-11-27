package shared

import (
	"regexp"
	"strings"

	"github.com/gosimple/unidecode"
)

// Slugify is a custom slug function to match current python implementation
func Slugify(value string) string {
	encoded := unidecode.Unidecode(value)

	reTrim := regexp.MustCompile(`[^\w\s-]`)
	reReplace := regexp.MustCompile(`[-\s]+`)

	trimmed := reTrim.ReplaceAllString(encoded, "")

	v := strings.ToLower(strings.TrimSpace(trimmed))
	return reReplace.ReplaceAllString(v, "_")
}

// EscapeChars escapes characters needed for HCL
func EscapeChars(value string) string {
	return strings.ReplaceAll(value, `\`, `\\`)
}
