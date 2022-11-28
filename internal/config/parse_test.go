package config

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/labd/mach-composer/internal/plugins"
	"github.com/labd/mach-composer/internal/utils"
	"github.com/labd/mach-composer/internal/variables"
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
          aws:
            account_id: 123456789
            region: eu-central-1
          components:
          - name: your-component
            variables:
              FOO_VAR: my-value
			  BAR_VAR: ${var.foo}
			  MULTIPLE_VARS: ${var.foo.bar} ${var.bar.foo}
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

	vars := variables.NewVariables()
	vars.Set("foo", "foobar")
	vars.Set("foo.bar", "1")
	vars.Set("bar.foo", "2")
	config, err := parseConfig(context.Background(), data, vars, "main.yml")
	if err != nil {
		t.Error(err)
	}

	expected := &MachConfig{
		MachComposer: MachComposer{
			Version: "1.0.0",
		},
		Global: Global{
			Environment:            "test",
			Cloud:                  "aws",
			TerraformStateProvider: "aws",
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
				Components: []SiteComponent{
					{
						Name: "your-component",
						Variables: map[string]any{
							"FOO_VAR":       "my-value",
							"BAR_VAR":       "foobar",
							"MULTIPLE_VARS": "1 2",
						},
						Secrets: map[string]any{
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
				Branch:       "",
				Integrations: []string{"aws", "commercetools"},
				Endpoints: map[string]string{
					"internal": "internal",
				},
			},
		},
		ExtraFiles: map[string][]byte{},
		Variables:  vars,
	}
	assert.Equal(t, expected.Global, config.Global)
	assert.Equal(t, expected.Sites, config.Sites)
	assert.Equal(t, expected.Components, config.Components)
	assert.Equal(t, expected.Variables, config.Variables)
}

func TestParseMissingVars(t *testing.T) {
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
          aws:
            account_id: 123456789
            region: eu-central-1
          components:
          - name: your-component
            variables:
              FOO_VAR: my-value
			  BAR_VAR: ${var.foo}
			  MULTIPLE_VARS: ${var.foo.bar} ${var.bar.foo}
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
    `))

	// Empty variables, it should fail because var.foo cannot be resolved
	vars := variables.Variables{}
	_, err := parseConfig(context.Background(), data, &vars, "main.yml")
	assert.Error(t, err)
}

func TestParseComponentsNodeInline(t *testing.T) {
	var intermediate struct {
		Components yaml.Node `yaml:"components"`
	}

	data := []byte(utils.TrimIndent(`
        components:
        - name: your-component
          source: "git::https://github.com/<username>/<your-component>.git//terraform"
          version: 0.1.0
  `))

	err := yaml.Unmarshal(data, &intermediate)
	require.NoError(t, err)

	cfg := &MachConfig{
		Plugins: plugins.NewPluginRepository(),
		Global: Global{
			Cloud: "aws",
		},
	}
	err = parseComponentsNode(cfg, &intermediate.Components, "main.yml")
	require.NoError(t, err)
	assert.Len(t, cfg.Components, 1)
	assert.Equal(t, "your-component", cfg.Components[0].Name)
}

func TestParseComponentsNodeInclude(t *testing.T) {
	utils.FS = afero.NewMemMapFs()
	utils.AFS = &afero.Afero{Fs: utils.FS}

	content := utils.TrimIndent(`
    - name: your-component
      source: "git::https://github.com/<username>/<your-component>.git//terraform"
      version: 0.1.0
	`)

	utils.AFS.WriteFile("components.yml", []byte(content), 0644)

	var intermediate struct {
		Components yaml.Node `yaml:"components"`
	}

	data := []byte(utils.TrimIndent(`
        components: ${include(components.yml)}
  `))

	err := yaml.Unmarshal(data, &intermediate)
	require.NoError(t, err)

	cfg := &MachConfig{
		Plugins: plugins.NewPluginRepository(),
		Global: Global{
			Cloud: "aws",
		},
	}
	err = parseComponentsNode(cfg, &intermediate.Components, "main.yml")
	require.NoError(t, err)
	assert.Len(t, cfg.Components, 1)
	assert.Equal(t, "your-component", cfg.Components[0].Name)
}
