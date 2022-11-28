package variables

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTemplateVariable(t *testing.T) {
	type test struct {
		input  string
		output string
	}

	tests := []test{
		{input: "${component.foo.endpoint}", output: "${module.foo.endpoint}"},
		{input: "foo ${component.foo.endpoint} bar", output: "foo ${module.foo.endpoint} bar"},
		{input: "  ${component.foo.endpoint} ", output: "${module.foo.endpoint}"},
		{input: "${component.foo.endpoint}${component.bar.other}", output: "${module.foo.endpoint}${module.bar.other}"},
	}

	for _, tc := range tests {
		value, err := parseComponentVariable(tc.input)
		assert.NoError(t, err)
		assert.Equal(t, tc.output, value)
	}
}

func TestInterpolateComponentVars(t *testing.T) {
	data := map[string]any{
		"my-simple-case": "${component.foobar.var}",
		"my-string-case": "-> ${component.foobar.var} <-",
		"my-nested-case": map[string]any{
			"my-simple-case": "${component.foobar.var}",
			"my-string-case": "-> ${component.foobar.var} - ${component.foobar.bar}<-",
		},
	}
	expected := map[string]any{
		"my-simple-case": "${module.foobar.var}",
		"my-string-case": "-> ${module.foobar.var} <-",
		"my-nested-case": map[string]any{
			"my-simple-case": "${module.foobar.var}",
			"my-string-case": "-> ${module.foobar.var} - ${module.foobar.bar}<-",
		},
	}
	result, err := InterpolateComponentVars(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
