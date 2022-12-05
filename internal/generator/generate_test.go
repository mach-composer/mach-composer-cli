package generator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/labd/mach-composer/internal/config"
	"github.com/labd/mach-composer/internal/plugins"
	"github.com/labd/mach-composer/internal/utils"
)

func TestRenderSite(t *testing.T) {
	data := []byte(utils.TrimIndent(`
	---
	mach_composer:
	  version: 1.0.0
	  plugins: {}
	global:
	  environment: test
	sites:
	- identifier: my-site
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
`))

	document := &yaml.Node{}
	err := yaml.Unmarshal(data, document)
	require.NoError(t, err)

	cfg, err := config.ParseConfig(
		context.Background(),
		document,
		config.ParseOptions{
			Filename: "main.yml",
			Plugins:  plugins.NewPluginRepository(),
		})
	if err != nil {
		t.Error(err)
	}

	body, err := renderSite(cfg, &cfg.Sites[0])
	assert.NoError(t, err)
	assert.NotEmpty(t, body)
}
