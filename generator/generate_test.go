package generator

import (
	"fmt"
	"testing"

	"github.com/labd/mach-composer-go/config"
	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {

	cfg := config.MachConfig{
		MachComposer: config.MachComposer{
			Version: "1.0.0",
		},
		Global: config.Global{
			Environment: "test",
			Cloud:       "aws",
			TerraformConfig: config.TerraformConfig{
				AwsRemoteState: &config.AWSTFState{
					Bucket:    "your bucket",
					KeyPrefix: "mach",
					Region:    "eu-central-1",
				},
			},
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
				Commercetools: &config.CommercetoolsSettings{
					ProjectKey:   "my-site",
					ClientID:     "<client-id>",
					ClientSecret: "<client-secret>",
					Scopes:       "manage_api_clients:my-site manage_project:my-site view_api_clients:my-site",
					ProjectSettings: &config.CommercetoolsProjectSettings{
						Languages:  []string{"en-GB", "nl-NL"},
						Currencies: []string{"GBP", "EUR"},
						Countries:  []string{"GB", "NL"},
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
				AWS: &config.SiteAWS{
					AccountID: "123456789",
					Region:    "eu-central-1",
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
	}

	config.Process(&cfg)

	body, err := Render(&cfg, &cfg.Sites[0])
	assert.NoError(t, err)

	fmt.Println(body)
}
