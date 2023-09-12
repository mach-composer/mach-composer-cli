package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTemplateVariable(t *testing.T) {
	type test struct {
		input  string
		output string
		ref    []string
	}

	tests := []test{
		{input: "${component.foo.endpoint}", output: "${module.foo.endpoint}", ref: []string{"foo"}},
		{input: "foo ${component.foo.endpoint} bar", output: "foo ${module.foo.endpoint} bar", ref: []string{"foo"}},
		{input: "  ${component.foo.endpoint} ", output: "${module.foo.endpoint}", ref: []string{"foo"}},
		{
			input:  "${component.foo.endpoint}${component.bar.other}",
			output: "${module.foo.endpoint}${module.bar.other}",
			ref:    []string{"foo", "bar"},
		},
	}

	for _, tc := range tests {
		value, ref, err := parseComponentVariable(tc.input)
		assert.NoError(t, err)
		assert.Equal(t, tc.output, value)
		assert.Equal(t, tc.ref, ref)
	}
}

func TestInterpolateComponentVarsSuccess(t *testing.T) {
	var testCases = []struct {
		Description string
		Input       map[string]any
		Expected    map[string]any
		Refs        []string
	}{
		{
			Description: "Simple test case",
			Input:       map[string]any{"key": "${component.foobar.var}"},
			Expected:    map[string]any{"key": "${module.foobar.var}"},
			Refs:        []string{"foobar"},
		},
		{
			Description: "String test case",
			Input:       map[string]any{"key": "-> ${component.foobar.var} <-"},
			Expected:    map[string]any{"key": "-> ${module.foobar.var} <-"},
			Refs:        []string{"foobar"},
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
				},
			},
			Expected: map[string]any{
				"key": map[string]any{
					"key-1": "${module.foobar.var}",
					"key-2": "-> ${module.foobar.var} - ${module.foobar.bar}<-",
				},
			},
			Refs: []string{"foobar"},
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
			Refs: []string{"site-topgeschenken"},
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
			Refs: []string{"site-topgeschenken"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			result, refs, err := interpolateComponentVar(testCase.Input)
			assert.NoError(t, err)
			assert.Equal(t, testCase.Expected, result)
			assert.Equal(t, testCase.Refs, refs)
		})
	}
}
