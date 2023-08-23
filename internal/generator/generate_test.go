package generator

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/plugins"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func TestRenderSite(t *testing.T) {
	content := []byte(utils.TrimIndent(`
	---
	mach_composer:
	  version: 1
	  plugins: {}
	global:
	  environment: test
	  terraform_config: {}
	  cloud: "aws"
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
	utils.FS = afero.NewMemMapFs()
	utils.AFS = &afero.Afero{Fs: utils.FS}

	err := utils.AFS.WriteFile("main.yml", []byte(content), 0644)
	require.NoError(t, err)

	pr := plugins.NewPluginRepository()
	err = pr.Add("aws", plugins.NewMockPlugin())
	require.NoError(t, err)

	cfg, err := config.Open(
		context.Background(), "main.yml", &config.ConfigOptions{
			Plugins: pr,
		})
	require.NoError(t, err)

	body, err := renderSite(context.TODO(), cfg, &cfg.Sites[0])
	assert.NoError(t, err)
	assert.NotEmpty(t, body)
}
