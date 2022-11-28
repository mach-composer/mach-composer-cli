package generator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/labd/mach-composer/internal/config"
	"github.com/labd/mach-composer/internal/utils"
)

func TestRenderSite(t *testing.T) {
	data := []byte(utils.TrimIndent(`
	---
	mach_composer:
	  version: 1.0.0
	global:
	  environment: test
	  cloud: aws
	  terraform_config:
		aws_remote_state:
	sites:
	- identifier: my-site
	  endpoints:
		main: api.my-site.nl
		internal:
		  url: internal-api.my-site.nl
		  throttling_burst_limit: 5000
		  throttling_rate_limit: 10000
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
	cfg, err := config.ParseConfig(context.Background(), data, nil, "main.yml")
	if err != nil {
		t.Error(err)
	}

	body, err := renderSite(cfg, &cfg.Sites[0])
	assert.NoError(t, err)
	assert.NotEmpty(t, body)
}
