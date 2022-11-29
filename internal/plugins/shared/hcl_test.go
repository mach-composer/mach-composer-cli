package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerializeToHCL(t *testing.T) {
	type test struct {
		input  any
		output string
	}

	tests := []test{
		{input: "foobar", output: `variables = "foobar"` + "\n"},
		{input: `c:\foo\bar`, output: `variables = "c:\\foo\\bar"` + "\n"},
		{input: 1, output: "variables = 1\n"},
		{input: 1.5, output: "variables = 1.5\n"},
		{input: true, output: "variables = true\n"},
		{input: false, output: "variables = false\n"},
		{input: []string{"${foo}", "${bar}"}, output: "variables = [foo, bar]\n"},

		{input: map[string]any{
			"my-key": map[string]any{
				"my-value": []string{
					"nl-NL",
					"en-US",
				},
			},
		},
			output: "variables = {\n  my-key = {\n    my-value = [\"nl-NL\", \"en-US\"]\n  }\n}\n"},
	}

	for _, tc := range tests {
		value := SerializeToHCL("variables", tc.input)
		assert.Equal(t, tc.output, value)
	}
}
