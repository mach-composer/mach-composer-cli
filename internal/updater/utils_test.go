package updater

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

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
	}

	utils.FS = afero.NewMemMapFs()
	utils.AFS = &afero.Afero{Fs: utils.FS}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.refContent != "" {
				err := utils.AFS.WriteFile("ref.yaml", []byte(tc.refContent), 0644)
				require.NoError(t, err)
			}

			_, filename, err := loadRefData(context.Background(), tc.node, "./")
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, filename, tc.refFilename)
				assert.NoError(t, err)

				expectedData := []byte{}
				require.NoError(t, tc.expected.Encode(expectedData))

				resultData := []byte{}
				require.NoError(t, tc.node.Encode(resultData))

				assert.Equal(t, expectedData, resultData)
			}
		})
	}
}
