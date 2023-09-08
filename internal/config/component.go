package config

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

var varComponentRegex = regexp.MustCompile(`\${(component(?:\.[^}]+)+)}`)

func interpolateComponentVar(value any) (any, []string, error) {
	var result any
	var refs []string

	switch v := value.(type) {
	case string:
		parsed, ref, err := parseComponentVariable(v)
		if err != nil {
			return nil, refs, err
		}
		result = parsed
		refs = append(refs, ref...)
	case []any:
		var slice []any
		for _, element := range v {
			ipElement, newRefs, err := interpolateComponentVar(element)
			if err != nil {
				return nil, refs, err
			}
			slice = append(slice, ipElement)
			refs = append(refs, slices.Compact(newRefs)...)
		}

		result = slice
	case map[string]any:
		var ipMap = map[string]any{}
		for key, element := range v {
			ipElement, newRefs, err := interpolateComponentVar(element)
			if err != nil {
				return nil, refs, err
			}
			ipMap[key] = ipElement
			refs = append(refs, slices.Compact(newRefs)...)
		}

		result = ipMap
	default:
		result = v
	}

	return result, refs, nil
}

// parseTemplateVariable replaces `${component.my-component.my-var}` with the
// terraform string `module.my-component.my-var` so that components can uses
// output of one component for another component.
func parseComponentVariable(raw string) (string, []string, error) {
	val := strings.TrimSpace(raw)
	var refs []string
	matches := varComponentRegex.FindAllStringSubmatch(val, 20)
	if len(matches) == 0 {
		result := val
		return result, refs, nil
	}

	for _, match := range matches {
		parts := strings.SplitN(match[1], ".", 3)
		if len(parts) < 3 {
			return "", refs, fmt.Errorf(
				"invalid variable '%s'; "+
					"When using a ${component...} variable it has to consist of 2 parts; "+
					"component-name.output-name",
				match[1])
		}

		replacement := fmt.Sprintf("${module.%s.%s}", parts[1], parts[2])
		val = strings.ReplaceAll(val, match[0], replacement)
		refs = append(refs, parts[1])
	}

	return val, refs, nil
}
