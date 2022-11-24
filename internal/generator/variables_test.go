package generator

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
		{input: "${component.foo.endpoint}", output: "module.foo.endpoint"},
		{input: "foo ${component.foo.endpoint} bar", output: `"foo ${module.foo.endpoint} bar"`},
		{input: "  ${component.foo.endpoint} ", output: "module.foo.endpoint"},
		{input: "${component.foo.endpoint}${component.bar.other}", output: `"${module.foo.endpoint}${module.bar.other}"`},
	}

	for _, tc := range tests {
		value, err := ParseTemplateVariable(tc.input)
		assert.NoError(t, err)
		assert.Equal(t, tc.output, value)
	}
}
