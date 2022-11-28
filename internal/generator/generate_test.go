package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/labd/mach-composer/internal/config"
	"github.com/labd/mach-composer/internal/plugins"
)

func TestRender(t *testing.T) {
	cfg := config.MachConfig{
		MachComposer: config.MachComposer{
			Version: "1.0.0",
		},
		Global: config.Global{
			Environment:            "test",
			Cloud:                  "aws",
			TerraformStateProvider: "aws",
		},
		Sites: []config.Site{
			{
				Name:       "",
				Identifier: "my-site",
				RawEndpoints: map[string]any{
					"main": "api.my-site.nl",
					"internal": map[string]any{
						"throttling_burst_limit": 5000,
						"throttling_rate_limit":  10000,
						"url":                    "internal-api.my-site.nl",
					},
				},
				Components: []config.SiteComponent{
					{
						Name: "your-component",
						Variables: map[string]any{
							"FOO_VAR": "my-value",
						},
						Secrets: map[string]any{
							"MY_SECRET": "secretvalue",
						},
					},
				},
			},
		},
		Components: []config.Component{
			{
				Name:         "your-component",
				Source:       "git::https://github.com/<username>/<your-component>.git//terraform",
				Version:      "0.1.0",
				Integrations: []string{"aws", "commercetools"},
			},
		},
		Plugins: plugins.NewPluginRepository(),
	}
	cfg.Plugins.Load("aws", "internal")

	config.ProcessConfig(&cfg)

	body, err := Render(&cfg, &cfg.Sites[0])
	assert.NoError(t, err)
	assert.NotEmpty(t, body)
}
