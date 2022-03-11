package config

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {

	data := []byte(dedent.Dedent(`
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
          integrations:
            - aws
            - commercetools
    `))

	fmt.Println(string(data))
	config, err := Parse(data)
	if err != nil {
		t.Error(err)
	}
	spew.Dump(config)

	expected := &Root{
		Global: Global{
			Environment: "test",
			Cloud:       "aws",
		},
		Sites: []Site{
			{
				Name:       "",
				Identifier: "my-site",
				CommercetoolsSettings: CommercetoolsSettings{
					ProjectKey:   "my-site",
					ClientID:     "<client-id>",
					ClientSecret: "<client-secret>",
					Scopes:       "manage_api_clients:my-site manage_project:my-site view_api_clients:my-site",
				},
				Components: []SiteComponent{
					{
						Name: "your-component",
						Variables: map[string]string{
							"FOO_VAR": "my-value",
						},
						Secrets: map[string]string{
							"MY_SECRET": "secretvalue",
						},
					},
				},
			},
		},
		Components: []Component{
			{
				Name:         "your-component",
				Source:       "git::https://github.com/<username>/<your-component>.git//terraform",
				Version:      "0.1.0",
				Integrations: []string{"aws", "commercetools"},
			},
		},
	}
	assert.Equal(t, expected, config)
}
