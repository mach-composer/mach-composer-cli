package variables

import (
	"fmt"
	"regexp"
	"strings"
)

var varComponentRegex = regexp.MustCompile(`\${(component(?:\.[^\}]+)+)}`)

func interpolateComponentVar(value any) (any, error) {
	var result any

	switch v := value.(type) {
	case string:
		parsed, err := parseComponentVariable(v)
		if err != nil {
			return nil, err
		}
		result = parsed
	case []any:
		var slice []any
		for _, element := range v {
			ipElement, err := interpolateComponentVar(element)
			if err != nil {
				return nil, err
			}
			slice = append(slice, ipElement)
		}

		result = slice
	case map[string]any:
		var ipMap = map[string]any{}
		for key, element := range v {
			ipElement, err := interpolateComponentVar(element)
			if err != nil {
				return nil, err
			}
			ipMap[key] = ipElement
		}

		result = ipMap
	default:
		result = v
	}

	return result, nil
}

func InterpolateComponentVars(data map[string]any) (map[string]any, error) {
	result := map[string]any{}

	for key, value := range data {
		ipVal, err := interpolateComponentVar(value)
		if err != nil {
			return nil, err
		}
		result[key] = ipVal
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
