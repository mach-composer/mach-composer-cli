package graph

import (
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
	"testing"
)

func TestSiteComponent_Hash_Ok(t *testing.T) {
	su := NewSiteComponent(nil, "", "", "", nil,
		config.SiteConfig{},
		config.SiteComponentConfig{Name: "b", Definition: &config.ComponentConfig{Name: "b", Source: "testdata/dirhash"}},
	)
	hash, err := su.Hash()

	assert.NoError(t, err)
	assert.Equal(t, "6c477d3042a9c36f088df7375cce8fed8e3a6c71d1c0da17d36e69593b8aafd7", hash, "Hashes should be equal")
}

func TestSiteComponent_HasChanges_HashNotFound(t *testing.T) {
	su := NewSiteComponent(nil, "", "", "", nil,
		config.SiteConfig{},
		config.SiteComponentConfig{Name: "b", Definition: &config.ComponentConfig{Name: "b", Source: "testdata/dirhash"}},
	)
	su.outputs = cty.ObjectVal(map[string]cty.Value{})

	changed, err := su.HasChanges()
	assert.NoError(t, err)
	assert.True(t, changed)
}

func TestSiteComponent_HasChanges_Error(t *testing.T) {
	su := NewSiteComponent(nil, "", "", "", nil,
		config.SiteConfig{},
		config.SiteComponentConfig{Name: "b", Definition: &config.ComponentConfig{Name: "b", Source: "testdata/dirhash"}},
	)
	su.outputs = cty.ObjectVal(map[string]cty.Value{
		"hash": cty.StringVal("some-hash"),
	})

	_, err := su.HasChanges()
	assert.Error(t, err)
}

func TestSiteComponent_HasChanges_True(t *testing.T) {
	su := NewSiteComponent(nil, "", "", "", nil,
		config.SiteConfig{},
		config.SiteComponentConfig{Name: "b", Definition: &config.ComponentConfig{Name: "b", Source: "testdata/dirhash"}},
	)
	su.outputs = cty.ObjectVal(map[string]cty.Value{
		"hash": cty.ObjectVal(map[string]cty.Value{
			"sensitive": cty.BoolVal(false),
			"value":     cty.StringVal("different-hash"),
			"type":      cty.StringVal("some-type"),
		}),
	})

	changed, err := su.HasChanges()
	assert.NoError(t, err)
	assert.True(t, changed)
}

func TestSiteComponent_HasChanges_False(t *testing.T) {
	su := NewSiteComponent(nil, "", "", "", nil,
		config.SiteConfig{},
		config.SiteComponentConfig{Name: "b", Definition: &config.ComponentConfig{Name: "b", Source: "testdata/dirhash"}},
	)
	su.outputs = cty.ObjectVal(map[string]cty.Value{
		"hash": cty.ObjectVal(map[string]cty.Value{
			"sensitive": cty.BoolVal(false),
			"value":     cty.StringVal("6c477d3042a9c36f088df7375cce8fed8e3a6c71d1c0da17d36e69593b8aafd7"),
			"type":      cty.StringVal("some-type"),
		}),
	})

	changed, err := su.HasChanges()
	assert.NoError(t, err)
	assert.False(t, changed)
}
