package config

import (
	"fmt"
	"testing"

	"github.com/labd/mach-composer-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {

	data := []byte(utils.TrimIndent(`
        ---
        mach_composer:
          version: 1.0.0
        global:
          environment: test
          terraform_config:
            aws_remote_state:
              bucket: "your bucket"
              key_prefix: mach
          cloud: aws
        sites:
        - identifier: my-site
          endpoints:
            main: api.my-site.nl
            internal:
              url: internal-api.my-site.nl
              throttling_burst_limit: 5000
              throttling_rate_limit: 10000
          commercetools:
            project_key: my-site
            client_id: "<client-id>"
            client_secret: "<client-secret>"
            scopes: manage_api_clients:my-site manage_project:my-site view_api_clients:my-site
            project_settings:
              languages:
                - en-GB
                - nl-NL
              currencies:
                - GBP
                - EUR
              countries:
                - GB
                - NL
          aws:
            account_id: 123456789
            region: eu-central-1
          components:
          - name: your-component
            variables:
              FOO_VAR: my-value
            secrets:
              MY_SECRET: secretvalue
        components:
        - name: your-component
          source: "git::https://github.com/<username>/<your-component>.git//terraform"
          version: 0.1.0
          endpoints:
            internal: internal
          integrations:
            - aws
            - commercetools
    `))

	fmt.Println(string(data))
	config, err := Parse(data)
	if err != nil {
		t.Error(err)
	}

	expected := &MachConfig{
		MachComposer: MachComposer{
			Version: "1.0.0",
		},
		Global: Global{
			Environment: "test",
			Cloud:       "aws",
			TerraformConfig: TerraformConfig{
				AwsRemoteState: &AWSTFState{
					Bucket:    "your bucket",
					KeyPrefix: "mach",
				},
			},
		},
		Sites: []Site{
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
				Commercetools: &CommercetoolsSettings{
					ProjectKey:   "my-site",
					ClientID:     "<client-id>",
					ClientSecret: "<client-secret>",
					Scopes:       "manage_api_clients:my-site manage_project:my-site view_api_clients:my-site",
					ProjectSettings: &CommercetoolsProjectSettings{
						Languages:  []string{"en-GB", "nl-NL"},
						Currencies: []string{"GBP", "EUR"},
						Countries:  []string{"GB", "NL"},
					},
				},
				Components: []SiteComponent{
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
				AWS: &SiteAWS{
					AccountID: "123456789",
					Region:    "eu-central-1",
				},
			},
		},
		Components: []Component{
			{
				Name:         "your-component",
				Source:       "git::https://github.com/<username>/<your-component>.git//terraform",
				Version:      "0.1.0",
				Integrations: []string{"aws", "commercetools"},
				Endpoints: map[string]string{
					"internal": "internal",
				},
			},
		},
	}
	assert.Equal(t, expected, config)
}
