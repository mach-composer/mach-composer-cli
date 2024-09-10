package config

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/mach-composer/mach-composer-cli/internal/config/variable"
	"testing"

	"github.com/mach-composer/mach-composer-cli/internal/plugins"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	config, err := Open(context.Background(), "testdata/configs/basic/main.yaml", &ConfigOptions{
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
		Filename: "main.yaml",
		MachComposer: MachComposer{
			Version: "1.0.0",
			Deployment: Deployment{
				Type: DeploymentSite,
			},
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
							"FOO_VAR":       variable.MustCreateNewScalarVariable("my-value"),
							"BAR_VAR":       variable.MustCreateNewScalarVariable("${var.foo}"),
							"MULTIPLE_VARS": variable.MustCreateNewScalarVariable("${var.foo.bar} ${var.bar.foo}"),
						},
						Secrets: variable.VariablesMap{
							"MY_SECRET": variable.MustCreateNewScalarVariable("secretvalue"),
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

func TestOpenComplex(t *testing.T) {
	pr := plugins.NewPluginRepository()
	err := pr.Add("my-plugin", plugins.NewPluginV1Adapter(plugins.NewMockPluginV1()))
	require.NoError(t, err)
	err = pr.Add("aws", plugins.NewPluginV1Adapter(plugins.NewMockPluginV1()))
	require.NoError(t, err)

	config, err := Open(context.Background(), "testdata/configs/complex/main.yaml", &ConfigOptions{
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
		Filename: "main.yaml",
		MachComposer: MachComposer{
			Version: "1.0.0",
			Plugins: map[string]MachPluginConfig{
				"aws":       {Source: "mach-composer/aws", Version: "0.1.0"},
				"my-plugin": {Source: "mach-composer/my-plugin", Version: "0.1.0"},
			},
			Deployment: Deployment{
				Type: DeploymentSite,
			},
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
							"FOO_VAR":       variable.MustCreateNewScalarVariable("my-value"),
							"BAR_VAR":       variable.MustCreateNewScalarVariable("${var.foo}"),
							"MULTIPLE_VARS": variable.MustCreateNewScalarVariable("${var.foo.bar} ${var.bar.foo}"),
						},
						Secrets: variable.VariablesMap{
							"MY_SECRET": variable.MustCreateNewScalarVariable("secretvalue"),
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

func TestOpenComponentIncludeFunc(t *testing.T) {
	config, err := Open(context.Background(), "testdata/configs/component_include_func/main.yaml", &ConfigOptions{
		Validate:      false,
		NoResolveVars: true,
		Plugins:       plugins.NewPluginRepository(),
	})
	require.NoError(t, err)

	assert.Len(t, config.Components, 1)
	assert.Equal(t, "your-component", config.Components[0].Name)
}

func TestOpenComponentRef(t *testing.T) {
	config, err := Open(context.Background(), "testdata/configs/component_ref/main.yaml", &ConfigOptions{
		Validate:      false,
		NoResolveVars: true,
		Plugins:       plugins.NewPluginRepository(),
	})
	require.NoError(t, err)

	assert.Len(t, config.Components, 1)
	assert.Equal(t, "your-component", config.Components[0].Name)
}

func TestAliased(t *testing.T) {
	config, err := Open(context.Background(), "testdata/configs/basic_alias/main.yaml", &ConfigOptions{
		Validate:      false,
		NoResolveVars: true,
		Plugins:       plugins.NewPluginRepository(),
	})
	require.NoError(t, err)

	assert.Len(t, config.Components, 2)
	assert.Equal(t, "your-component", config.Components[0].Name)
	assert.Equal(t, "0.1.0", config.Components[0].Version)
	assert.Equal(t, "your-component-aliased", config.Components[1].Name)
	assert.Equal(t, "0.1.0", config.Components[1].Version)
}

func TestComponentRefAliased(t *testing.T) {
	config, err := Open(context.Background(), "testdata/configs/component_ref_alias/main.yaml", &ConfigOptions{
		Validate:      false,
		NoResolveVars: true,
		Plugins:       plugins.NewPluginRepository(),
	})
	require.NoError(t, err)

	assert.Len(t, config.Components, 2)
	assert.Equal(t, "your-component", config.Components[0].Name)
	assert.Equal(t, "0.1.0", config.Components[0].Version)
	assert.Equal(t, "your-component-aliased", config.Components[1].Name)
	assert.Equal(t, "0.1.0", config.Components[1].Version)
}
