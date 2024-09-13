package variable

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeVariablesMaps(t *testing.T) {
	tests := []struct {
		name     string
		maps     []VariablesMap
		expected VariablesMap
	}{
		{
			name: "merge two maps with different keys",
			maps: []VariablesMap{
				{"key1": MustCreateNewScalarVariable("value1")},
				{"key2": MustCreateNewScalarVariable("value2")},
			},
			expected: VariablesMap{
				"key1": MustCreateNewScalarVariable("value1"),
				"key2": MustCreateNewScalarVariable("value2"),
			},
		},
		{
			name: "merge two maps with same keys",
			maps: []VariablesMap{
				{"key1": MustCreateNewScalarVariable("value1")},
				{"key1": MustCreateNewScalarVariable("value2")},
			},
			expected: VariablesMap{
				"key1": MustCreateNewScalarVariable("value2"),
			},
		},
		{
			name: "merge three maps with same keys",
			maps: []VariablesMap{
				{"key1": MustCreateNewScalarVariable("value1")},
				{"key2": MustCreateNewScalarVariable("value2")},
				{"key1": MustCreateNewScalarVariable("value3")},
			},
			expected: VariablesMap{
				"key1": MustCreateNewScalarVariable("value3"),
				"key2": MustCreateNewScalarVariable("value2"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := MergeVariablesMaps(tc.maps...)
			assert.Equal(t, tc.expected, result)
		})
	}
}
