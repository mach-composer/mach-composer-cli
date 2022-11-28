package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	plugin.SetSiteEndpointsConfig("my-site", data)

	result := plugin.endpointsConfigs["my-site"]["internal"]
	assert.Equal(t, intRef(5000), result.ThrottlingBurstLimit)
	assert.Equal(t, intRef(10000), result.ThrottlingRateLimit)

}
