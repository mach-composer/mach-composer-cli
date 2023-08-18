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

func TestInterpolateComponentVarsSuccess(t *testing.T) {
	var testCases = []struct {
		Description string
		Input       map[string]any
		Expected    map[string]any
	}{
		{
			Description: "Simple test case",
			Input:       map[string]any{"key": "${component.foobar.var}"},
			Expected:    map[string]any{"key": "${module.foobar.var}"},
		},
		{
			Description: "String test case",
			Input:       map[string]any{"key": "-> ${component.foobar.var} <-"},
			Expected:    map[string]any{"key": "-> ${module.foobar.var} <-"},
		},
		{
			Description: "Integer test case",
			Input:       map[string]any{"key": 123},
			Expected:    map[string]any{"key": 123},
		},
		{
			Description: "Boolean test case",
			Input:       map[string]any{"key": true},
			Expected:    map[string]any{"key": true},
		},
		{
			Description: "Float test case",
			Input:       map[string]any{"key": 1.1},
			Expected:    map[string]any{"key": 1.1},
		},
		{
			Description: "Nested map test case",
			Input: map[string]any{
				"key": map[string]any{
					"key-1": "${component.foobar.var}",
					"key-2": "-> ${component.foobar.var} - ${component.foobar.bar}<-",
				}},
			Expected: map[string]any{
				"key": map[string]any{
					"key-1": "${module.foobar.var}",
					"key-2": "-> ${module.foobar.var} - ${module.foobar.bar}<-",
				}},
		},
		{
			Description: "Nested slice test case",
			Input: map[string]any{
				"key": []interface{}{
					`${component.site-topgeschenken.service_link_map["topbloemen-nl"]}`,
					`${component.site-topgeschenken.service_link_map["topbloemen-be"]}`,
				},
			},
			Expected: map[string]any{
				"key": []interface{}{
					`${module.site-topgeschenken.service_link_map["topbloemen-nl"]}`,
					`${module.site-topgeschenken.service_link_map["topbloemen-be"]}`,
				}},
		},
		{
			Description: "Slice with nested map test case",
			Input: map[string]any{
				"key": []interface{}{
					map[string]any{
						"key": `${component.site-topgeschenken.service_link_map["topbloemen-nl"]}`,
					},
					map[string]any{
						"key": `${component.site-topgeschenken.service_link_map["topbloemen-be"]}`,
					},
				},
			},
			Expected: map[string]any{
				"key": []interface{}{
					map[string]any{"key": `${module.site-topgeschenken.service_link_map["topbloemen-nl"]}`},
					map[string]any{"key": `${module.site-topgeschenken.service_link_map["topbloemen-be"]}`},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			result, err := InterpolateComponentVars(testCase.Input)
			assert.NoError(t, err)
			assert.Equal(t, testCase.Expected, result)
		})
	}
}
