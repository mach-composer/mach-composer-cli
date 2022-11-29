package shared

import (
	"fmt"
	"regexp"
)

var tfProviderRegex = regexp.MustCompile(`([!=<>~]*)(.*)`)

func VersionConstraint(expression string) string {
	// match = TF_PROVIDER_RE.match(value or default_version)
	// operator, version = match.groups()
	// return f"{operator or '~>'} {version}"
	var operator string

	val := expression
	match := tfProviderRegex.FindStringSubmatch(expression)
	if match != nil {
		operator = match[1]
		val = match[2]
	}

	if len(operator) == 0 {
		operator = "~>"
	}

	return fmt.Sprintf("%s %s", operator, val)
}
