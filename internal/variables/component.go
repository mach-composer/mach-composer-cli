package variables

import (
	"fmt"
	"regexp"
	"strings"
)

var varComponentRegex = regexp.MustCompile(`\${(component(?:\.[^\}]+)+)}`)

func InterpolateComponentVars(data map[string]any) (map[string]any, error) {
	result := map[string]any{}

	for key, value := range data {
		switch v := value.(type) {
		case string:
			parsed, err := parseComponentVariable(v)
			if err != nil {
				return nil, err
			}
			result[key] = parsed
		case map[string]any:
			parsed, err := InterpolateComponentVars(v)
			if err != nil {
				return nil, err
			}
			result[key] = parsed
		default:
			result[key] = value
		}
	}
	return result, nil
}

// parseTemplateVariable replaces `${component.my-component.my-var}` with the
// terraform string `module.my-component.my-var` so that components can uses
// output of one component for another component.
func parseComponentVariable(raw string) (string, error) {
	val := strings.TrimSpace(raw)
	matches := varComponentRegex.FindAllStringSubmatch(val, 20)
	if len(matches) == 0 {
		result := val
		return result, nil
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

	return val, nil
}
