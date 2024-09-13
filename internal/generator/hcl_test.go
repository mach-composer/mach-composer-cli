package generator

import (
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/config/variable"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerializeToHCL(t *testing.T) {
	type test struct {
		input          variable.VariablesMap
		output         string
		deploymentType config.DeploymentType
	}

	tests := []test{
		{
			input: variable.VariablesMap{
				"foo": variable.MustCreateNewScalarVariable("bar"),
			},
			output:         "variables = {\n  foo = \"bar\"\n}\n",
			deploymentType: config.DeploymentSite,
		},
		{
			input: variable.VariablesMap{
				"foo": variable.MustCreateNewScalarVariable("c:\\foo\\bar"),
			},
			output:         "variables = {\n  foo = \"c:\\\\foo\\\\bar\"\n}\n",
			deploymentType: config.DeploymentSite,
		},
		{
			input: variable.VariablesMap{
				"foo": variable.MustCreateNewScalarVariable(1),
			},
			output:         "variables = {\n  foo = 1\n}\n",
			deploymentType: config.DeploymentSite,
		},
		{
			input: variable.VariablesMap{
				"foo": variable.MustCreateNewScalarVariable(1.5),
			},
			output:         "variables = {\n  foo = 1.5\n}\n",
			deploymentType: config.DeploymentSite,
		},
		{
			input: variable.VariablesMap{
				"foo": variable.MustCreateNewScalarVariable(true),
			},
			output:         "variables = {\n  foo = true\n}\n",
			deploymentType: config.DeploymentSite,
		},
		{
			input: variable.VariablesMap{
				"foo": variable.MustCreateNewScalarVariable(false),
			},
			output:         "variables = {\n  foo = false\n}\n",
			deploymentType: config.DeploymentSite,
		},
		{
			input: variable.VariablesMap{
				"foo": variable.NewSliceVariable([]variable.Variable{
					variable.MustCreateNewScalarVariable("${foo}"),
					variable.MustCreateNewScalarVariable("${bar}"),
				}),
			},
			output:         "variables = {\n  foo = [foo, bar]\n}\n",
			deploymentType: config.DeploymentSite,
		},
		{
			input: variable.VariablesMap{
				"foo": variable.MustCreateNewScalarVariable("${data.sops_external.variables.data[\"foo-bar\"]}"),
			},
			output:         "variables = {\n  foo = data.sops_external.variables.data[\"foo-bar\"]\n}\n",
			deploymentType: config.DeploymentSite,
		},
		{
			input: variable.VariablesMap{
				"foo": variable.NewMapVariable(map[string]variable.Variable{
					"63000012": variable.MustCreateNewScalarVariable("foobar"),
				}),
			},
			output:         "variables = {\n  foo = {\n    \"63000012\" = \"foobar\"\n  }\n}\n",
			deploymentType: config.DeploymentSite,
		},
		{
			input: variable.VariablesMap{
				"foo": variable.NewMapVariable(map[string]variable.Variable{
					"my-key": variable.NewMapVariable(map[string]variable.Variable{
						"my-value": variable.NewSliceVariable([]variable.Variable{
							variable.MustCreateNewScalarVariable("nl-NL"),
							variable.MustCreateNewScalarVariable("en-US"),
						}),
					}),
				}),
			},
			output:         "variables = {\n  foo = {\n    my-key = {\n      my-value = [\"nl-NL\", \"en-US\"]\n    }\n  }\n}\n",
			deploymentType: config.DeploymentSite,
		},
	}

	for _, tc := range tests {
		value, err := serializeToHCL("variables", tc.input, tc.deploymentType, nil, "test-1")
		assert.NoError(t, err)
		assert.Equal(t, tc.output, value)
	}
}
