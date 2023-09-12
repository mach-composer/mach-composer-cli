package generator

import (
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerializeToHCL(t *testing.T) {
	type test struct {
		input  config.SiteComponentVars
		output string
	}

	tests := []test{
		{
			input: config.SiteComponentVars{
				"foo": &config.SiteComponentVar{Interpolated: "bar"},
			},
			output: "variables = {\n  foo = \"bar\"\n}\n",
		},
		{
			input: config.SiteComponentVars{
				"foo": &config.SiteComponentVar{Interpolated: "c:\\foo\\bar"},
			},
			output: "variables = {\n  foo = \"c:\\\\foo\\\\bar\"\n}\n",
		},
		{
			input: config.SiteComponentVars{
				"foo": &config.SiteComponentVar{Interpolated: 1},
			},
			output: "variables = {\n  foo = 1\n}\n",
		},
		{
			input: config.SiteComponentVars{
				"foo": &config.SiteComponentVar{Interpolated: 1.5},
			},
			output: "variables = {\n  foo = 1.5\n}\n",
		},
		{
			input: config.SiteComponentVars{
				"foo": &config.SiteComponentVar{Interpolated: true},
			},
			output: "variables = {\n  foo = true\n}\n",
		},
		{
			input: config.SiteComponentVars{
				"foo": &config.SiteComponentVar{Interpolated: false},
			},
			output: "variables = {\n  foo = false\n}\n",
		},
		{
			input: config.SiteComponentVars{
				"foo": &config.SiteComponentVar{Interpolated: []string{"${foo}", "${bar}"}},
			},
			output: "variables = {\n  foo = [foo, bar]\n}\n",
		},
		{
			input: config.SiteComponentVars{
				"foo": &config.SiteComponentVar{Interpolated: `${data.sops_external.variables.data["foo-bar"]}`},
			},
			output: "variables = {\n  foo = data.sops_external.variables.data[\"foo-bar\"]\n}\n",
		},
		{
			input: config.SiteComponentVars{
				"foo": &config.SiteComponentVar{Interpolated: map[int]any{63000012: "foobar"}},
			},
			output: "variables = {\n  foo = {\n    \"63000012\" = \"foobar\"\n  }\n}\n",
		},
		{
			input: config.SiteComponentVars{
				"foo": &config.SiteComponentVar{
					Interpolated: map[string]any{
						"my-key": map[string]any{
							"my-value": []string{
								"nl-NL", "en-US",
							},
						},
					},
				},
			},
			output: "variables = {\n  foo = {\n    my-key = {\n      my-value = [\"nl-NL\", \"en-US\"]\n    }\n  }\n}\n",
		},
	}

	for _, tc := range tests {
		value, err := serializeToHCL("variables", tc.input)
		assert.NoError(t, err)
		assert.Equal(t, tc.output, value)
	}
}
