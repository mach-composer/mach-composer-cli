package graph

import (
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/config/variable"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashSiteComponentConfigOk(t *testing.T) {
	val, _ := variable.NewScalarVariable("value1")

	h, err := hashSiteComponentConfig(config.SiteComponentConfig{
		Name: "site-component-1",
		Variables: variable.VariablesMap{
			"var1": val,
		},
		Secrets: variable.VariablesMap{},
		Definition: &config.ComponentConfig{
			Name:   "site-component-1",
			Source: "testdata/dirhash",
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "930c819f53a9b36fa7dc3857acb2d567dc3dd0c8c3f41922e957ea9fe4b066ea", h)
}

func TestHashSiteComponentConfigChanged(t *testing.T) {
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

	h1, err := hashSiteComponentConfig(cfg)

	cfg.DependsOn = []string{"site-component-2"}

	h2, err := hashSiteComponentConfig(cfg)

	assert.NoError(t, err)
	assert.NotEqual(t, h1, h2)
}

func TestHashSiteComponentConfigGithubSource(t *testing.T) {
	val, _ := variable.NewScalarVariable("value1")
	cfg := config.SiteComponentConfig{
		Name: "site-component-1",
		Variables: variable.VariablesMap{
			"var1": val,
		},
		Secrets: variable.VariablesMap{},
		Definition: &config.ComponentConfig{
			Name:   "site-component-1",
			Source: "git@github.com:hashicorp/example.git",
		},
	}

	h, err := hashSiteComponentConfig(cfg)

	assert.NoError(t, err)
	assert.Equal(t, "5fa49450032642f53ca6df3cd853530cbd0f2b6468d250e7603980007b91bf7a", h)
}
