package generator

import (
	"testing"

	"github.com/flosch/pongo2/v5"
	"github.com/stretchr/testify/assert"
)

func TestFilterTFValue(t *testing.T) {
	type test struct {
		input  *pongo2.Value
		output *pongo2.Value
	}

	tests := []test{
		{input: pongo2.AsValue("foobar"), output: pongo2.AsSafeValue(`"foobar"`)},
		{input: pongo2.AsValue(1), output: pongo2.AsSafeValue("1")},
		{input: pongo2.AsValue(1.5), output: pongo2.AsSafeValue("1.5")},
		{input: pongo2.AsValue(true), output: pongo2.AsSafeValue("true")},
		{input: pongo2.AsValue(false), output: pongo2.AsSafeValue("false")},
		{input: pongo2.AsValue([]string{"foo", "bar"}), output: pongo2.AsSafeValue(`["foo", "bar"]`)},
	}

	for _, tc := range tests {
		value, err := filterTFValue(tc.input, nil)
		assert.Nil(t, err)
		assert.True(t, tc.output.EqualValueTo(value))
	}
}
