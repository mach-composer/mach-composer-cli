package generator

import (
	"fmt"
	"regexp"
	"strings"
)

var varRegex = regexp.MustCompile(`\${(component(?:\.[^\}]+)+)}`)

func ParseTemplateVariable(raw string) (string, error) {
	val := strings.TrimSpace(raw)
	org := val

	matches := varRegex.FindAllStringSubmatch(val, 20)
	if len(matches) == 0 {
		return val, nil
	}

	for _, match := range matches {
		parts := strings.SplitN(match[1], ".", 3)
		if len(parts) < 3 {
			return "", fmt.Errorf(
				"invalid variable '%s'; "+
					"When using a ${component...} variable it has to consist of 2 parts; "+
					"component-name.output-name",
				match[1])
		}

		replacement := fmt.Sprintf("${module.%s.%s}", parts[1], parts[2])
		val = strings.ReplaceAll(val, match[0], replacement)
	}

	var result string
	if len(matches) == 1 && len(matches[0][0]) == len(org) {
		result = val[2 : len(val)-1]
	} else {
		result = fmt.Sprintf(`"%s"`, val)
	}

	return result, nil
}
