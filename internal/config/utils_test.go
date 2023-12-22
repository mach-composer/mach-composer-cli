package config

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/mach-composer/mach-composer-cli/internal/utils"
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
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := mapYamlNodes(tc.nodes)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestLoadRefData(t *testing.T) {
	tests := []struct {
		name        string
		node        *yaml.Node
		refContent  string
		refFilename string
		expected    *yaml.Node
		wantErr     bool
	}{
		{
			name: "mapping node with $ref",
			node: &yaml.Node{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{Value: "$ref"},
					{Value: "ref.yaml"},
				},
			},
			refContent:  "key: value",
			refFilename: "ref.yaml",
			expected: &yaml.Node{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{Value: "key"},
					{Value: "value"},
				},
			},
			wantErr: false,
		},
		{
			name: "mapping node with ${include()}",
			node: &yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: "${include(ref.yaml)}",
			},
			refContent:  "key: value",
			refFilename: "ref.yaml",
			expected: &yaml.Node{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{Value: "key"},
					{Value: "value"},
				},
			},
			wantErr: false,
		},
		{
			name: "mapping node with nested $ref",
			node: &yaml.Node{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{Value: "$ref"},
					{Value: "ref.yaml"},
				},
			},
			refContent: utils.TrimIndent(`
				some-node:
					other-node:
						key: value
			`),
			refFilename: "ref.yaml",
			expected: &yaml.Node{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{Value: "key"},
					{Value: "value"},
				},
			},
			wantErr: false,
		},
		{
			name: "mapping node without $ref",
			node: &yaml.Node{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{Value: "key"},
					{Value: "value"},
				},
			},
			refContent:  "",
			refFilename: "",
			expected: &yaml.Node{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{Value: "key"},
					{Value: "value"},
				},
			},
			wantErr: false,
		},
		{
			name: "error loading ref document",
			node: &yaml.Node{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{Value: "$ref"},
					{Value: "other.yaml"},
				},
			},
			refContent:  "",
			refFilename: "other.yaml",
			expected:    nil,
			wantErr:     true,
		},
		{
			name: "error loading ref document with include",
			node: &yaml.Node{
				Kind: yaml.ScalarNode,
				Content: []*yaml.Node{
					{Value: "${include(ref.yaml)}"},
				},
			},
			refContent: "key: value",
			wantErr:    true,
		},
	}

	utils.FS = afero.NewMemMapFs()
	utils.AFS = &afero.Afero{Fs: utils.FS}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.refContent != "" {
				err := utils.AFS.WriteFile("ref.yaml", []byte(tc.refContent), 0644)
				require.NoError(t, err)
			}

			_, filename, err := LoadRefData(context.Background(), tc.node, "./")
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, filename, tc.refFilename)
				assert.NoError(t, err)

				var expectedData []byte
				require.NoError(t, tc.expected.Encode(expectedData))

				var resultData []byte
				require.NoError(t, tc.node.Encode(resultData))

				assert.Equal(t, expectedData, resultData)
			}
		})
	}
}
