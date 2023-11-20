package variable

import (
	"github.com/mach-composer/mach-composer-cli/internal/state"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewScalarVariable(t *testing.T) {
	type test struct {
		input string
		ref   []string
	}

	tests := []test{
		{input: "${component.foo.endpoint}", ref: []string{"foo"}},
		{input: "foo ${component.foo.endpoint} bar", ref: []string{"foo"}},
		{input: "  ${component.foo.endpoint} ", ref: []string{"foo"}},
		{input: "${component.foo.endpoint}${component.bar.other}", ref: []string{"foo", "bar"}},
	}

	for _, tc := range tests {
		value, err := NewScalarVariable(tc.input)
		assert.NoError(t, err)
		assert.Equal(t, tc.ref, value.references)
	}
}

func TestModuleTransformFunc(t *testing.T) {
	type test struct {
		input  string
		output string
	}

	tests := []test{
		{input: "${component.foo.endpoint}", output: "${module.foo.endpoint}"},
		{input: "foo ${component.foo.endpoint} bar", output: "foo ${module.foo.endpoint} bar"},
		{input: "  ${component.foo.endpoint} ", output: "${module.foo.endpoint}"},
		{
			input:  "${component.foo.endpoint}${component.bar.other}",
			output: "${module.foo.endpoint}${module.bar.other}",
		},
	}

	for _, tc := range tests {
		value, err := NewScalarVariable(tc.input)
		assert.NoError(t, err)

		res, err := value.TransformValue(ModuleTransformFunc())
		assert.NoError(t, err)

		assert.Equal(t, tc.output, res)
	}
}

func TestRemoteStateTransformFunc(t *testing.T) {
	type test struct {
		input  string
		output string
	}

	tests := []test{
		{
			input:  "${component.foo.endpoint}",
			output: "${data.terraform_remote_state.my_state.outputs.foo.variables.endpoint}"},
		{
			input:  "foo ${component.foo.endpoint} bar",
			output: "foo ${data.terraform_remote_state.my_state.outputs.foo.variables.endpoint} bar",
		},
		{
			input:  "  ${component.foo.endpoint} ",
			output: "${data.terraform_remote_state.my_state.outputs.foo.variables.endpoint}",
		},
		{
			input:  "${data.terraform_remote_state.my_state.outputs.foo.endpoint}${data.terraform_remote_state.my_state.outputs.bar.other}",
			output: "${data.terraform_remote_state.my_state.outputs.foo.endpoint}${data.terraform_remote_state.my_state.outputs.bar.other}",
		},
	}

	r := state.NewRepository()
	err := r.Add("my_state", nil)
	r.Alias("my_state", "foo")
	assert.NoError(t, err)

	for _, tc := range tests {
		value, err := NewScalarVariable(tc.input)
		assert.NoError(t, err)

		res, err := value.TransformValue(RemoteStateTransformFunc(r))
		assert.NoError(t, err)

		assert.Equal(t, tc.output, res)
	}
}

func TestInterpolateComponentVarsSuccess(t *testing.T) {
	var testCases = []struct {
		Description string
		Input       *VariablesMap
		Expected    map[string]any
		Refs        []string
	}{
		{
			Description: "Simple test case",
			Input: &VariablesMap{
				"key": MustCreateNewScalarVariable(t, "${component.foobar.var}"),
			},
			Expected: map[string]any{"key": "${module.foobar.var}"},
			Refs:     []string{"foobar"},
		},
		{
			Description: "Scalar test case",
			Input: &VariablesMap{
				"key": MustCreateNewScalarVariable(t, "-> ${component.foobar.var} <-"),
			},
			Expected: map[string]any{"key": "-> ${module.foobar.var} <-"},
			Refs:     []string{"foobar"},
		},
		{
			Description: "Integer test case",
			Input: &VariablesMap{
				"key": MustCreateNewScalarVariable(t, 123),
			},
			Expected: map[string]any{"key": 123},
		},
		{
			Description: "Boolean test case",
			Input: &VariablesMap{
				"key": MustCreateNewScalarVariable(t, true),
			},
			Expected: map[string]any{"key": true},
		},
		{
			Description: "Float test case",
			Input: &VariablesMap{
				"key": MustCreateNewScalarVariable(t, 1.1),
			},
			Expected: map[string]any{"key": 1.1},
		},
		{
			Description: "Nested map test case",
			Input: &VariablesMap{
				"key": &MapVariable{
					Elements: map[string]Variable{
						"key-1": MustCreateNewScalarVariable(t, "${component.foobar.var}"),
						"key-2": MustCreateNewScalarVariable(t, "-> ${component.foobar.var} - ${component.foobar.bar}<-"),
					},
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
			Input: &VariablesMap{
				"key": &SliceVariable{
					Elements: []Variable{
						MustCreateNewScalarVariable(t, `${component.site-topgeschenken.service_link_map["topbloemen-nl"]}`),
						MustCreateNewScalarVariable(t, `${component.site-topgeschenken.service_link_map["topbloemen-be"]}`),
					},
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
			Input: &VariablesMap{
				"key": &SliceVariable{
					Elements: []Variable{
						&MapVariable{
							Elements: map[string]Variable{
								"key": MustCreateNewScalarVariable(t, `${component.site-topgeschenken.service_link_map["topbloemen-nl"]}`),
							},
						},
						&MapVariable{
							Elements: map[string]Variable{
								"key": MustCreateNewScalarVariable(t, `${component.site-topgeschenken.service_link_map["topbloemen-be"]}`),
							},
						},
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
			result, err := testCase.Input.Transform(ModuleTransformFunc())
			assert.NoError(t, err)
			assert.Equal(t, testCase.Expected, result)
			assert.Equal(t, testCase.Refs, testCase.Input.ListReferencedComponents())
		})
	}
}
