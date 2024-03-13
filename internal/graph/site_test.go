package graph

import (
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
	"testing"
)

func TestSite_Hash_NestedComponentConfigSorted(t *testing.T) {
	su := NewSite(nil, "", "", "", nil, config.SiteConfig{})
	su.NestedNodes = []config.SiteComponentConfig{
		{Name: "b", Definition: &config.ComponentConfig{Name: "b", Source: "testdata/dirhash"}},
		{Name: "a", Definition: &config.ComponentConfig{Name: "a", Source: "testdata/dirhash"}},
	}

	unsortedHash, err := su.Hash()

	s := NewSite(nil, "", "", "", nil, config.SiteConfig{})
	s.NestedNodes = []config.SiteComponentConfig{
		{Name: "a", Definition: &config.ComponentConfig{Name: "a", Source: "testdata/dirhash"}},
		{Name: "b", Definition: &config.ComponentConfig{Name: "b", Source: "testdata/dirhash"}},
	}
	sortedHash, err := s.Hash()

	assert.NoError(t, err)
	assert.Equal(t, unsortedHash, sortedHash, "Hashes should be equal")
}

func TestSite_HasChanges_NoHash(t *testing.T) {
	s := NewSite(nil, "", "", "", nil, config.SiteConfig{})
	s.NestedNodes = []config.SiteComponentConfig{
		{Name: "b", Definition: &config.ComponentConfig{Name: "b", Source: "testdata/dirhash"}},
		{Name: "a", Definition: &config.ComponentConfig{Name: "a", Source: "testdata/dirhash"}},
	}
	s.outputs = cty.ObjectVal(map[string]cty.Value{})

	changed, err := s.HasChanges()
	assert.NoError(t, err)
	assert.True(t, changed)
}

func TestSite_HasChanges_Error(t *testing.T) {
	s := NewSite(nil, "", "", "", nil, config.SiteConfig{})
	s.outputs = cty.ObjectVal(map[string]cty.Value{
		"hash": cty.StringVal("some-hash"),
	})

	_, err := s.HasChanges()
	assert.Error(t, err)
}

func TestSite_HasChanges_True(t *testing.T) {
	s := NewSite(nil, "", "", "", nil, config.SiteConfig{})
	s.NestedNodes = []config.SiteComponentConfig{
		{Name: "b", Definition: &config.ComponentConfig{Name: "b", Source: "testdata/dirhash"}},
		{Name: "a", Definition: &config.ComponentConfig{Name: "a", Source: "testdata/dirhash"}},
	}
	s.outputs = cty.ObjectVal(map[string]cty.Value{
		"hash": cty.ObjectVal(map[string]cty.Value{
			"sensitive": cty.BoolVal(false),
			"value":     cty.StringVal("different-hash"),
			"type":      cty.StringVal("some-type"),
		}),
	})

	changed, err := s.HasChanges()
	assert.NoError(t, err)
	assert.True(t, changed)
}

func TestSite_HasChanges_False(t *testing.T) {
	s := NewSite(nil, "", "", "", nil, config.SiteConfig{})
	s.NestedNodes = []config.SiteComponentConfig{
		{Name: "b", Definition: &config.ComponentConfig{Name: "b", Source: "testdata/dirhash"}},
		{Name: "a", Definition: &config.ComponentConfig{Name: "a", Source: "testdata/dirhash"}},
	}

	s.outputs = cty.ObjectVal(map[string]cty.Value{
		"hash": cty.ObjectVal(map[string]cty.Value{
			"sensitive": cty.BoolVal(false),
			"value":     cty.StringVal("6b70e89e46522af5c0dad35f403ed12c13e891cef0f1360024be3645754fa53e"),
			"type":      cty.StringVal("some-type"),
		}),
	})

	changed, err := s.HasChanges()
	assert.NoError(t, err)
	assert.False(t, changed)
}
