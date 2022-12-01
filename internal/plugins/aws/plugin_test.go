package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func intRef(val int) *int {
	return &val
}

func TestSetSiteEndpointsConfig(t *testing.T) {
	data := map[string]any{
		"internal": map[string]any{
			"url": "example.org",
			"aws": map[string]any{
				"throttling_burst_limit": 5000,
				"throttling_rate_limit":  10000,
			},
		},
	}

	plugin := NewAWSPlugin()
	err := plugin.SetSiteConfig("my-site", map[string]any{})
	require.NoError(t, err)

	err = plugin.SetComponentConfig("my-component", map[string]any{
		"integrations": []string{"aws"},
	})
	require.NoError(t, err)

	err = plugin.SetComponentEndpointsConfig("my-component", map[string]string{
		"internal": "internal",
	})
	require.NoError(t, err)

	err = plugin.SetSiteEndpointsConfig("my-site", data)
	require.NoError(t, err)

	result, err := plugin.RenderTerraformResources("my-site")
	require.NoError(t, err)
	assert.Contains(t, result, "throttling_burst_limit = 5000")
}
