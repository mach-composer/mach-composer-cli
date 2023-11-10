package config

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

var ignoreOpts = []cmp.Option{
	cmpopts.IgnoreUnexported(MachConfig{}, Variables{}, variable.ScalarVariable{}),
	cmpopts.IgnoreFields(MachConfig{}, "StateRepository", "Plugins", "Variables"),
}

func TestOpenBasic(t *testing.T) {
	config, err := Open(context.Background(), "testdata/configs/basic.yaml", &ConfigOptions{
		Validate:      false,
		NoResolveVars: true,
		Plugins:       plugins.NewPluginRepository(),
	})
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
		Filename: "basic.yaml",
		MachComposer: MachComposer{
			Version:    "1.0.0",
			Deployment: &Deployment{Type: DeploymentSite},
		},
		Global: GlobalConfig{
			Environment: "test",
		},
		Sites: []SiteConfig{
			{
				Name:       "",
				Deployment: &Deployment{Type: DeploymentSite},
				Identifier: "my-site",
				Components: []SiteComponentConfig{
					{
						Name:       "your-component",
						Deployment: &Deployment{Type: DeploymentSite},
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
	}

	assert.True(t, cmp.Equal(config, expected, ignoreOpts...), cmp.Diff(config, expected, ignoreOpts...))
}

func TestParseComplex(t *testing.T) {
	pr := plugins.NewPluginRepository()
	err := pr.Add("my-plugin", plugins.NewMockPlugin())
	require.NoError(t, err)
	err = pr.Add("aws", plugins.NewMockPlugin())
	require.NoError(t, err)

	config, err := Open(context.Background(), "testdata/configs/complex.yaml", &ConfigOptions{
		Validate:      false,
		NoResolveVars: true,
		Plugins:       pr,
	})
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
		Filename: "complex.yaml",
		MachComposer: MachComposer{
			Version: "1.0.0",
			Plugins: map[string]MachPluginConfig{
				"aws":       {Source: "mach-composer/aws", Version: "0.1.0"},
				"my-plugin": {Source: "mach-composer/my-plugin", Version: "0.1.0"},
			},
			Deployment: &Deployment{Type: DeploymentSite},
		},
		Global: GlobalConfig{
			Environment:            "test",
			Cloud:                  "aws",
			TerraformStateProvider: "aws",
			TerraformConfig: &TerraformConfig{
				RemoteState: map[string]any{"plugin": "aws"},
				Providers:   nil,
			},
		},
		Sites: []SiteConfig{
			{
				Name:       "",
				Deployment: &Deployment{Type: DeploymentSite},
				Identifier: "my-site",
				RawEndpoints: map[string]any{
					"main": "api.my-site.nl",
					"internal": map[string]any{
						"aws": map[string]any{
							"throttling_burst_limit": 5000,
							"throttling_rate_limit":  10000,
						},
						"url": "internal-api.my-site.nl",
					},
				},
				Components: []SiteComponentConfig{
					{
						Name:       "your-component",
						Deployment: &Deployment{Type: DeploymentSite},
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
	}

	assert.True(t, cmp.Equal(config, expected, ignoreOpts...), cmp.Diff(config, expected, ignoreOpts...))
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
