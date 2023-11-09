package variable

import (
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/state"
	"regexp"
	"strings"
)

var varComponentRegex = regexp.MustCompile(`\${(component(?:\.[^}]+)+)}`)

type ScalarVariable struct {
	baseVariable
	content    string
	references []string
}

func (v *ScalarVariable) TransformValue(f TransformValueFunc) (any, error) {
	return f(v.content)
}

func (v *ScalarVariable) Content() string {
	return v.content
}

func (v *ScalarVariable) ReferencedComponents() []string {
	return v.references
}

func parseValue(v string) ([]string, error) {
	val := strings.TrimSpace(v)
	matches := varComponentRegex.FindAllStringSubmatch(val, 20)
	if len(matches) == 0 {
		return nil, nil
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple references found in variable '%s'", val)
	}

	return strings.SplitN(matches[0][1], ".", 3), nil
}

func parseReference(v string) (string, bool, error) {
	val, err := parseValue(v)
	if err != nil {
		return "", false, err
	}
	if val == nil {
		return v, false, nil
	}
	return val[1], true, nil
}

func ModuleTransformFunc() TransformValueFunc {
	return func(value string) (any, error) {
		parts, err := parseValue(value)
		if err != nil {
			return nil, err
		}

		if len(parts) == 0 {
			return value, nil
		}
		replacement := fmt.Sprintf("${module.%s.%s}", parts[1], parts[2])
		return replacement, nil
	}
}

func RemoteStateTransformFunc(repository *state.Repository) TransformValueFunc {
	return func(value string) (any, error) {
		parts, err := parseValue(value)
		if err != nil {
			return nil, err
		}

		if len(parts) == 0 {
			return value, nil
		}

		var stateKey, ok = repository.Key(parts[1])
		if !ok {
			return nil, fmt.Errorf("state key '%s' not found", parts[1])
		}

		replacement := fmt.Sprintf(`${data.terraform_remote_state.%s.outputs.%s.%s}`, stateKey, parts[1], parts[2])
		return replacement, nil
	}
}
