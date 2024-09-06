package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestNodeAsMap(t *testing.T) {
	tests := []struct {
		name     string
		node     *yaml.Node
		expected map[string]any
		wantErr  bool
	}{
		{
			name: "valid mapping node",
			node: &yaml.Node{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Value: "key1"},
					{Kind: yaml.ScalarNode, Value: "value1"},
					{Kind: yaml.ScalarNode, Value: "key2"},
					{Kind: yaml.ScalarNode, Value: "123"},
				},
			},
			expected: map[string]any{
				"key1": "value1",
				"key2": 123,
			},
		},
		{
			name: "non-mapping node",
			node: &yaml.Node{
				Kind: yaml.ScalarNode,
				Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Value: "key1"},
					{Kind: yaml.ScalarNode, Value: "value1"},
				},
			},
			expected: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := nodeAsMap(tc.node)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMapYamlNodes(t *testing.T) {
	tests := []struct {
		name     string
		nodes    []*yaml.Node
		expected map[string]*yaml.Node
	}{
		{
			name: "even number of nodes",
			nodes: []*yaml.Node{
				{Value: "key1"},
				{Value: "value1"},
				{Value: "key2"},
				{Value: "value2"},
			},
			expected: map[string]*yaml.Node{
				"key1": {Value: "value1"},
				"key2": {Value: "value2"},
			},
		},
		{
			name: "odd number of nodes",
			nodes: []*yaml.Node{
				{Value: "key1"},
				{Value: "value1"},
				{Value: "key2"},
			},
			expected: nil,
		},
		{
			name: "anchor node",
			nodes: []*yaml.Node{
				{Value: "key1"},
				{Value: "value1"},
				{Value: "<<"},
				{
					Alias: &yaml.Node{
						Content: []*yaml.Node{
							{Value: "key2"},
							{Value: "value2"},
						}},
				},
			},
			expected: map[string]*yaml.Node{
				"key1": {Value: "value1"},
				"key2": {Value: "value2"},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := MapYamlNodes(tc.nodes)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestLoadIncludeDocument(t *testing.T) {
	tests := []struct {
		name        string
		node        *yaml.Node
		cwd         string
		expectedErr bool
	}{
		{
			name: "valid include syntax with valid components.yaml",
			node: &yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: "${include(components.yaml)}",
			},
			cwd:         "testdata",
			expectedErr: false,
		},
		{
			name: "valid include syntax with valid components.yml",
			node: &yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: "${include(components.yml)}",
			},
			cwd:         "testdata",
			expectedErr: false,
		},
		{
			name: "invalid include syntax",
			node: &yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: "${include(invalid-syntax}",
			},
			cwd:         "testdata",
			expectedErr: true,
		},
		{
			name: "valid include syntax with non-existent file",
			node: &yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: "${include(nonexistent.yml)}",
			},
			cwd:         "testdata",
			expectedErr: true,
		},
		{
			name: "valid include syntax with invalid YAML file",
			node: &yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: "${include(invalid.yml)}",
			},
			cwd:         "testdata",
			expectedErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := LoadIncludeDocument(tc.node, tc.cwd)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
