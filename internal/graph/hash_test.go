package graph

import (
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/config/variable"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashSiteComponentConfigOk(t *testing.T) {
	val, _ := variable.NewScalarVariable("value1")

	h, err := HashSiteComponent(&SiteComponent{
		SiteComponentConfig: config.SiteComponentConfig{
			Name: "site-component-1",
			Variables: variable.VariablesMap{
				"var1": val,
			},
			Secrets: variable.VariablesMap{},
			Definition: &config.ComponentConfig{
				Name:   "site-component-1",
				Source: "testdata/dirhash",
			},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "fbdafc38bdec19d852c8a31e3e09a63df076fb1ab7a82b4e920429346bf8652f", h)
}

func TestHashSiteComponentConfigChanged(t *testing.T) {
	val, _ := variable.NewScalarVariable("value1")
	n := &SiteComponent{
		SiteComponentConfig: config.SiteComponentConfig{
			Name: "site-component-1",
			Variables: variable.VariablesMap{
				"var1": val,
			},
			Secrets: variable.VariablesMap{},
			Definition: &config.ComponentConfig{
				Name:   "site-component-1",
				Source: "testdata/dirhash",
			},
		},
	}

	h1, err := HashSiteComponent(n)

	n.SiteComponentConfig.DependsOn = []string{"site-component-2"}

	h2, err := HashSiteComponent(n)

	assert.NoError(t, err)
	assert.NotEqual(t, h1, h2)
}

func TestHashSiteComponentConfigGithubSource(t *testing.T) {
	val, _ := variable.NewScalarVariable("value1")
	n := &SiteComponent{
		SiteComponentConfig: config.SiteComponentConfig{
			Name: "site-component-1",
			Variables: variable.VariablesMap{
				"var1": val,
			},
			Secrets: variable.VariablesMap{},
			Definition: &config.ComponentConfig{
				Name:   "site-component-1",
				Source: "testdata/dirhash",
			},
		},
	}

	h, err := HashSiteComponent(n)

	assert.NoError(t, err)
	assert.Equal(t, "fbdafc38bdec19d852c8a31e3e09a63df076fb1ab7a82b4e920429346bf8652f", h)
}

func TestHashSiteComponentConfigVariableFileOK(t *testing.T) {
	val, _ := variable.NewScalarVariable("value1")

	h, err := HashSiteComponent(&SiteComponent{
		ProjectConfig: config.MachConfig{
			MachComposer: config.MachComposer{
				VariablesFile: "testdata/variables-1.yaml",
			},
		},
		SiteComponentConfig: config.SiteComponentConfig{
			Name: "site-component-1",
			Variables: variable.VariablesMap{
				"var1": val,
			},
			Secrets: variable.VariablesMap{},
			Definition: &config.ComponentConfig{
				Name:   "site-component-1",
				Source: "testdata/dirhash",
			},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "93270be1128022f9bad163959e8ff036f3b443991d8fbe463f51a9b01e46be11", h)
}

func TestHashSiteComponentConfigVariablesFileChanged(t *testing.T) {
	val, _ := variable.NewScalarVariable("value1")
	cfg := config.SiteComponentConfig{
		Name: "site-component-1",
		Variables: variable.VariablesMap{
			"var1": val,
		},
		Secrets: variable.VariablesMap{},
		Definition: &config.ComponentConfig{
			Name:   "site-component-1",
			Source: "testdata/dirhash",
		},
	}

	h1, err := HashSiteComponent(&SiteComponent{
		ProjectConfig: config.MachConfig{
			MachComposer: config.MachComposer{
				VariablesFile: "testdata/variables-1.yaml",
			},
		},
		SiteComponentConfig: cfg,
	})

	h2, err := HashSiteComponent(&SiteComponent{
		ProjectConfig: config.MachConfig{
			MachComposer: config.MachComposer{
				VariablesFile: "testdata/variables-2.yaml",
			},
		},
		SiteComponentConfig: cfg,
	})

	assert.NoError(t, err)
	assert.NotEqual(t, h1, h2)
}
