package config

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/config/variable"
	"testing"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/mach-composer/mach-composer-cli/internal/plugins"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	m.Run()
}

func TestParseBasic(t *testing.T) {
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
			  BAR_VAR: ${var.foo}
			  MULTIPLE_VARS: ${var.foo.bar} ${var.bar.foo}
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

	// Decode the yaml in an intermediate config file
	intermediate, err := newRawConfig("test.yml", document)
	require.NoError(t, err)

	intermediate.MachComposer.Deployment = &Deployment{
		Type: DeploymentSite,
	}

	vars := NewVariables()
	vars.Set("foo", "foobar")
	vars.Set("foo.bar", "1")
	vars.Set("bar.foo", "2")
	intermediate.variables = vars
	intermediate.plugins = plugins.NewPluginRepository()

	err = resolveVariables(context.Background(), intermediate, ".")
	require.NoError(t, err)

	config, err := resolveConfig(context.Background(), intermediate)
	require.NoError(t, err)

	component := ComponentConfig{
		Name:         "your-component",
		Source:       "git::https://github.com/<username>/<your-component>.git//terraform",
		Version:      "0.1.0",
		Branch:       "",
		Integrations: []string{},
		Endpoints:    nil,
	}

	expected := &MachConfig{
		MachComposer: MachComposer{
			Version: "1.0.0",
		},
		Sites: []SiteConfig{
			{
				Name:       "",
				Identifier: "my-site",
				Components: []SiteComponentConfig{
					{
						Name: "your-component",
						Variables: variable.VariablesMap{
							"FOO_VAR":       variable.MustCreateNewScalarVariable(t, "my-value"),
							"BAR_VAR":       variable.MustCreateNewScalarVariable(t, "foobar"),
							"MULTIPLE_VARS": variable.MustCreateNewScalarVariable(t, "1 2"),
						},
						Secrets: variable.VariablesMap{
							"MY_SECRET": variable.MustCreateNewScalarVariable(t, "my-secretvalue"),
						},
						Definition: &component,
					},
				},
			},
		},
		Components: []ComponentConfig{component},
		extraFiles: map[string][]byte{},
		Variables:  vars,
	}
	assert.Equal(t, expected.Global, config.Global)
	assert.Equal(t, expected.Sites, config.Sites)
	assert.Equal(t, expected.Components, config.Components)
	assert.Equal(t, expected.Variables, config.Variables)
}

func TestParse(t *testing.T) {
	data := []byte(utils.TrimIndent(`
        ---
        mach_composer:
          version: 1.0.0
        global:
          environment: test
          terraform_config:
			remote_state:
				plugin: "my-plugin"
			providers:
				aws: 3.0.0
          cloud: my-plugin
        sites:
        - identifier: my-site
          endpoints:
            main: api.my-site.nl
            internal:
              url: internal-api.my-site.nl
              throttling_burst_limit: 5000
              throttling_rate_limit: 10000
          my-plugin:
            some-key: 123456789
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
            - my-plugin
    `))

	vars := NewVariables()
	vars.Set("foo", "foobar")
	vars.Set("foo.bar", "1")
	vars.Set("bar.foo", "2")

	pluginRepo := plugins.NewPluginRepository()
	err := pluginRepo.Add("my-plugin", plugins.NewMockPlugin())
	require.NoError(t, err)

	document := &yaml.Node{}
	err = yaml.Unmarshal(data, document)
	require.NoError(t, err)

	intermediate, err := newRawConfig("main.yml", document)
	require.NoError(t, err)

	intermediate.plugins = pluginRepo
	intermediate.variables = vars

	err = resolveVariables(context.Background(), intermediate, ".")
	require.NoError(t, err)

	config, err := resolveConfig(context.Background(), intermediate)
	require.NoError(t, err)

	component := ComponentConfig{
		Name:         "your-component",
		Source:       "git::https://github.com/<username>/<your-component>.git//terraform",
		Version:      "0.1.0",
		Branch:       "",
		Integrations: []string{"my-plugin"},
		Endpoints: map[string]string{
			"internal": "internal",
		},
	}

	expected := &MachConfig{
		MachComposer: MachComposer{
			Version: "1.0.0",
		},
		Global: GlobalConfig{
			Environment:            "test",
			Cloud:                  "my-plugin",
			TerraformStateProvider: "my-plugin",
			TerraformConfig: &TerraformConfig{
				RemoteState: map[string]any{
					"plugin": "my-plugin",
				},
				Providers: map[string]string{
					"aws": "3.0.0",
				},
			},
		},
		Sites: []SiteConfig{
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
				Components: []SiteComponentConfig{
					{
						Name: "your-component",
						Variables: variable.VariablesMap{
							"FOO_VAR":       variable.MustCreateNewScalarVariable(t, "my-value"),
							"BAR_VAR":       variable.MustCreateNewScalarVariable(t, "foobar"),
							"MULTIPLE_VARS": variable.MustCreateNewScalarVariable(t, "1 2"),
						},
						Secrets: variable.VariablesMap{
							"MY_SECRET": variable.MustCreateNewScalarVariable(t, "my-secretvalue"),
						},
						Definition: &component,
					},
				},
			},
		},
		Components: []ComponentConfig{component},
		extraFiles: map[string][]byte{},
		Variables:  vars,
	}
	assert.Equal(t, expected.Global, config.Global)
	assert.Equal(t, expected.Sites, config.Sites)
	assert.Equal(t, expected.Components, config.Components)
	assert.Equal(t, expected.Variables, config.Variables)
}

func TestResolveMissingVar(t *testing.T) {
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

	document := &yaml.Node{}
	err := yaml.Unmarshal(data, document)
	require.NoError(t, err)

	intermediate, err := newRawConfig("main.yml", document)
	require.NoError(t, err)
	intermediate.variables = &Variables{}

	err = resolveVariables(context.Background(), intermediate, ".")
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
		Global: GlobalConfig{
			Cloud: "my-cloud",
		},
	}
	err = cfg.Plugins.Add("my-cloud", plugins.NewMockPlugin())
	require.NoError(t, err)

	err = parseComponentsNode(cfg, &intermediate.Components)
	require.NoError(t, err)
	assert.Len(t, cfg.Components, 1)
	assert.Equal(t, "your-component", cfg.Components[0].Name)
}

func TestParseComponentsNodeRef(t *testing.T) {
	utils.FS = afero.NewMemMapFs()
	utils.AFS = &afero.Afero{Fs: utils.FS}

	content := utils.TrimIndent(`
	components:
		- name: your-component
		  source: "git::https://github.com/<username>/<your-component>.git//terraform"
		  version: 0.1.0
	`)

	err := utils.AFS.WriteFile("components.yml", []byte(content), 0644)
	require.NoError(t, err)

	var intermediate struct {
		Components yaml.Node `yaml:"components"`
	}

	data := []byte(utils.TrimIndent(`
        components:
			$ref: components.yml#/components
  `))

	err = yaml.Unmarshal(data, &intermediate)
	require.NoError(t, err)

	componentNode, fileName, err := LoadRefData(context.Background(), &intermediate.Components, "")
	assert.NoError(t, err)
	assert.Equal(t, "components.yml", fileName)
	intermediate.Components = *componentNode

	cfg := &MachConfig{
		Plugins: plugins.NewPluginRepository(),
		Global: GlobalConfig{
			Cloud: "my-cloud",
		},
	}
	err = cfg.Plugins.Add("my-cloud", plugins.NewMockPlugin())
	require.NoError(t, err)

	err = parseComponentsNode(cfg, &intermediate.Components)
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

	err := utils.AFS.WriteFile("components.yml", []byte(content), 0644)
	require.NoError(t, err)

	var intermediate struct {
		Components yaml.Node `yaml:"components"`
	}

	data := []byte(utils.TrimIndent(`
        components: ${include(components.yml)}
  `))

	err = yaml.Unmarshal(data, &intermediate)
	require.NoError(t, err)

	componentNode, fileName, err := LoadRefData(context.Background(), &intermediate.Components, "")
	assert.NoError(t, err)
	assert.Equal(t, "components.yml", fileName)
	intermediate.Components = *componentNode

	cfg := &MachConfig{
		Plugins: plugins.NewPluginRepository(),
		Global: GlobalConfig{
			Cloud: "my-cloud",
		},
	}
	err = cfg.Plugins.Add("my-cloud", plugins.NewMockPlugin())
	require.NoError(t, err)

	err = parseComponentsNode(cfg, &intermediate.Components)
	require.NoError(t, err)
	assert.Len(t, cfg.Components, 1)
	assert.Equal(t, "your-component", cfg.Components[0].Name)
}
