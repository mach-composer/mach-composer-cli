package dependency

import (
	"github.com/dominikbraun/graph"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/config/variable"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromConfigSimple(t *testing.T) {
	cfg := &config.MachConfig{
		Filename: "main.yml",
		Sites: []config.SiteConfig{
			{
				Name:       "site 1",
				Identifier: "site-1",
				Components: []config.SiteComponentConfig{
					{
						Name: "site-component-1",
					},
				},
			},
		},
	}

	g, err := ToDependencyGraph(cfg)
	assert.NoError(t, err)

	o, _ := g.Order()
	assert.Equal(t, 3, o)

	e, _ := g.Edges()
	assert.Len(t, e, 2)

	site, err := g.Vertex("main.yml/site-1")
	assert.NoError(t, err)
	assert.Equal(t, "main.yml/site-1", site.Path())
	assert.IsType(t, &Site{}, site)
	assert.Equal(t, &cfg.Sites[0], site.(*Site).SiteConfig)

	sc, err := g.Vertex("main.yml/site-1/site-component-1")
	assert.NoError(t, err)
	assert.Equal(t, "main.yml/site-1/site-component-1", sc.Path())
	assert.IsType(t, &SiteComponent{}, sc)
	assert.Equal(t, &cfg.Sites[0].Components[0], sc.(*SiteComponent).SiteComponentConfig)

	am, _ := g.AdjacencyMap()
	assert.Equal(t, map[string]map[string]graph.Edge[string]{
		"main.yml": {
			"main.yml/site-1": graph.Edge[string]{
				Source:     "main.yml",
				Target:     "main.yml/site-1",
				Properties: graph.EdgeProperties{Attributes: map[string]string{}, Weight: 0, Data: interface{}(nil)},
			},
		},
		"main.yml/site-1": {
			"main.yml/site-1/site-component-1": graph.Edge[string]{
				Source:     "main.yml/site-1",
				Target:     "main.yml/site-1/site-component-1",
				Properties: graph.EdgeProperties{Attributes: map[string]string{}, Weight: 0, Data: interface{}(nil)},
			},
		},
		"main.yml/site-1/site-component-1": {}}, am)
}

func TestFromConfigExplicit(t *testing.T) {
	cfg := &config.MachConfig{
		Filename: "main.yml",
		Sites: []config.SiteConfig{
			{
				Name:       "site 1",
				Identifier: "site-1",
				Components: []config.SiteComponentConfig{
					{
						Name: "site-component-1",
					},
					{
						Name: "site-component-2",
						DependsOn: []string{
							"site-component-1",
						},
					},
				},
			},
		},
	}

	g, err := ToDependencyGraph(cfg)
	assert.NoError(t, err)

	am, _ := g.AdjacencyMap()
	assert.Equal(t, map[string]map[string]graph.Edge[string]{
		"main.yml": {
			"main.yml/site-1": graph.Edge[string]{
				Source:     "main.yml",
				Target:     "main.yml/site-1",
				Properties: graph.EdgeProperties{Attributes: map[string]string{}, Weight: 0, Data: interface{}(nil)},
			},
		},
		"main.yml/site-1": {
			"main.yml/site-1/site-component-1": graph.Edge[string]{
				Source:     "main.yml/site-1",
				Target:     "main.yml/site-1/site-component-1",
				Properties: graph.EdgeProperties{Attributes: map[string]string{}, Weight: 0, Data: interface{}(nil)},
			},
		},
		"main.yml/site-1/site-component-1": {
			"main.yml/site-1/site-component-2": graph.Edge[string]{
				Source:     "main.yml/site-1/site-component-1",
				Target:     "main.yml/site-1/site-component-2",
				Properties: graph.EdgeProperties{Attributes: map[string]string{}, Weight: 0, Data: interface{}(nil)},
			},
		},
		"main.yml/site-1/site-component-2": {},
	}, am)
}

func TestFromConfigInferred(t *testing.T) {
	cfg := &config.MachConfig{
		Filename: "main.yml",
		Sites: []config.SiteConfig{
			{
				Name:       "site 1",
				Identifier: "site-1",
				Components: []config.SiteComponentConfig{
					{
						Name: "site-component-1",
					},
					{
						Name: "site-component-2",
						Variables: variable.VariablesMap{
							"value-1": variable.MustCreateNewScalarVariable(t, ""),
						},
					},
				},
			},
		},
	}

	g, err := ToDependencyGraph(cfg)
	assert.NoError(t, err)

	am, _ := g.AdjacencyMap()
	assert.Equal(t, map[string]map[string]graph.Edge[string]{
		"main.yml": {
			"main.yml/site-1": graph.Edge[string]{
				Source:     "main.yml",
				Target:     "main.yml/site-1",
				Properties: graph.EdgeProperties{Attributes: map[string]string{}, Weight: 0, Data: interface{}(nil)},
			},
		},
		"main.yml/site-1": {
			"main.yml/site-1/site-component-1": graph.Edge[string]{
				Source:     "main.yml/site-1",
				Target:     "main.yml/site-1/site-component-1",
				Properties: graph.EdgeProperties{Attributes: map[string]string{}, Weight: 0, Data: interface{}(nil)},
			},
		},
		"main.yml/site-1/site-component-1": {
			"main.yml/site-1/site-component-2": graph.Edge[string]{
				Source:     "main.yml/site-1/site-component-1",
				Target:     "main.yml/site-1/site-component-2",
				Properties: graph.EdgeProperties{Attributes: map[string]string{}, Weight: 0, Data: interface{}(nil)},
			},
		},
		"main.yml/site-1/site-component-2": {},
	}, am)
}

func TestFromConfigCycleErr(t *testing.T) {
	cfg := &config.MachConfig{
		Filename: "main.yml",
		Sites: []config.SiteConfig{
			{
				Name:       "site 1",
				Identifier: "site-1",
				Components: []config.SiteComponentConfig{
					{
						Name: "site-component-1",
						DependsOn: []string{
							"site-component-2",
						},
					},
					{
						Name: "site-component-2",
						DependsOn: []string{
							"site-component-1",
						},
					},
				},
			},
		},
	}

	_, err := ToDependencyGraph(cfg)
	assert.Error(t, err)
	assert.IsType(t, &ValidationError{}, err)
	assert.Len(t, err.(*ValidationError).Errors, 1)
	assert.Equal(t,
		"failed to add dependency from main.yml/site-1/site-component-1 to main.yml/site-1/site-component-2: "+
			"edge would create a cycle",
		err.(*ValidationError).Errors[0].Error(),
	)
}
