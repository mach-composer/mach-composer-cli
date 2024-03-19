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
	assert.Equal(t, "28ea13834f4d36a8799267a372d3a8d4bb8353097fcbd3b2445005b702e9ac8b", h)
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
			Source: "git::github.com:hashicorp/example.git",
		},
	}

	h, err := hashSiteComponentConfig(cfg)

	assert.NoError(t, err)
	assert.Equal(t, "de87afc8419dcd29e3e8cbe2e47b5026593ac0975555fe3d0f341eb3e0cf5785", h)
}
