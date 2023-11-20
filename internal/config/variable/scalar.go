package variable

import (
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/state"
	"github.com/stretchr/testify/require"
	"regexp"
	"strings"
	"testing"
)

var varComponentRegex = regexp.MustCompile(`\${(component(?:\.[^}]+)+)}`)

type ScalarVariable struct {
	baseVariable
	Content    any
	references []string
}

func NewScalarVariable(content any) (*ScalarVariable, error) {
	var references []string
	if s, ok := content.(string); ok {
		parsedReferences, err := parseReferences(s)
		if err != nil {
			return nil, err
		}

		references = append(references, parsedReferences...)
	}

	return &ScalarVariable{baseVariable: baseVariable{typ: Scalar}, Content: content, references: references}, nil
}

func (v *ScalarVariable) TransformValue(f TransformValueFunc) (any, error) {
	return f(v.Content)
}

func (v *ScalarVariable) ReferencedComponents() []string {
	return v.references
}

func parseValues(v string) ([][]string, error) {
	val := strings.TrimSpace(v)
	matches := varComponentRegex.FindAllStringSubmatch(val, 20)
	if len(matches) == 0 {
		return nil, nil
	}

	var parsedValues [][]string
	for _, match := range matches {
		parts := strings.SplitN(match[1], ".", 3)
		if len(parts) < 3 {
			return nil, fmt.Errorf(
				"invalid variable '%s'; "+
					"When using a ${component...} variable it has to consist of 2 parts; "+
					"component-name.output-name",
				match[1])
		}

		parsedValues = append(parsedValues, parts)
	}

	return parsedValues, nil
}

func parseReferences(v string) ([]string, error) {
	val, err := parseValues(v)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}

	var references []string
	for _, v := range val {
		references = append(references, v[1])
	}

	return references, nil
}

func ModuleTransformFunc() TransformValueFunc {
	return func(value any) (any, error) {
		val, ok := value.(string)
		if !ok {
			return value, nil
		}

		parts, err := parseValues(val)
		if err != nil {
			return nil, err
		}

		if len(parts) == 0 {
			return value, nil
		}

		for _, part := range parts {
			replacement := fmt.Sprintf("module.%s.%s", part[1], part[2])
			val = strings.ReplaceAll(val, strings.Join(part, "."), replacement)
		}

		return strings.TrimSpace(val), nil
	}
}

func RemoteStateTransformFunc(repository *state.Repository) TransformValueFunc {
	return func(value any) (any, error) {
		val, ok := value.(string)
		if !ok {
			return value, nil
		}

		parts, err := parseValues(val)
		if err != nil {
			return nil, err
		}

		if len(parts) == 0 {
			return value, nil
		}

		for _, part := range parts {
			stateKey, exists := repository.Key(part[1])
			if !exists {
				return nil, fmt.Errorf("state key '%s' not found", part[1])
			}

			replacement := fmt.Sprintf(`data.terraform_remote_state.%s.outputs.%s.variables.%s`, stateKey, part[1], part[2])
			val = strings.ReplaceAll(val, strings.Join(part, "."), replacement)
		}
		return strings.TrimSpace(val), nil
	}
}

func MustCreateNewScalarVariable(t *testing.T, value any) *ScalarVariable {
	v, err := NewScalarVariable(value)
	require.NoError(t, err)
	return v
}
