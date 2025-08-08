package graph

import (
	"github.com/dominikbraun/graph"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/config/variable"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToDependencyGraphSimple(t *testing.T) {
	cfg := &config.MachConfig{
		Filename: "main",
		MachComposer: config.MachComposer{
			Deployment: config.Deployment{
				Type: config.DeploymentSiteComponent,
			},
		},

		Sites: []config.SiteConfig{
			{
				Name: "site 1",
				Deployment: &config.Deployment{
					Type: config.DeploymentSiteComponent,
				},
				Identifier: "site-1",
				Components: []config.SiteComponentConfig{
					{
						Name: "site-component-1",
						Deployment: &config.Deployment{
							Type: config.DeploymentSiteComponent,
						},
					},
				},
			},
		},
	}

	g, err := ToDependencyGraph(cfg, "")
	assert.NoError(t, err)

	o, _ := g.Order()
	assert.Equal(t, 3, o)

	e, _ := g.Edges()
	assert.Len(t, e, 2)

	site, err := g.Vertex("main/site-1")
	assert.NoError(t, err)
	assert.Equal(t, "main/site-1", site.Path())
	assert.IsType(t, &Site{}, site)
	assert.Equal(t, cfg.Sites[0], site.(*Site).SiteConfig)

	sc, err := g.Vertex("main/site-1/site-component-1")
	assert.NoError(t, err)
	assert.Equal(t, "main/site-1/site-component-1", sc.Path())
	assert.IsType(t, &SiteComponent{}, sc)
	assert.Equal(t, cfg.Sites[0].Components[0], sc.(*SiteComponent).SiteComponentConfig)

	am, _ := g.AdjacencyMap()
	assert.Equal(t, map[string]map[string]graph.Edge[string]{
		"main": {
			"main/site-1": graph.Edge[string]{
				Source:     "main",
				Target:     "main/site-1",
				Properties: graph.EdgeProperties{Attributes: map[string]string{}, Weight: 0, Data: interface{}(nil)},
			},
		},
		"main/site-1": {
			"main/site-1/site-component-1": graph.Edge[string]{
				Source:     "main/site-1",
				Target:     "main/site-1/site-component-1",
				Properties: graph.EdgeProperties{Attributes: map[string]string{}, Weight: 0, Data: interface{}(nil)},
			},
		},
		"main/site-1/site-component-1": {}}, am)
}

func TestToDependencyGraphExplicit(t *testing.T) {
	cfg := &config.MachConfig{
		Filename: "main",
		MachComposer: config.MachComposer{
			Deployment: config.Deployment{
				Type: config.DeploymentSiteComponent,
			},
		},
		Sites: []config.SiteConfig{
			{
				Name: "site 1",
				Deployment: &config.Deployment{
					Type: config.DeploymentSiteComponent,
				},
				Identifier: "site-1",
				Components: []config.SiteComponentConfig{
					{
						Name: "site-component-1",
						Deployment: &config.Deployment{
							Type: config.DeploymentSiteComponent,
						},
					},
					{
						Name: "site-component-2",
						Deployment: &config.Deployment{
							Type: config.DeploymentSiteComponent,
						},
						DependsOnKeys: []string{
							"site-component-1",
						},
					},
				},
			},
		},
	}

	g, err := ToDependencyGraph(cfg, "")
	assert.NoError(t, err)

	am, _ := g.AdjacencyMap()

	assert.Equal(t, map[string]map[string]graph.Edge[string]{
		"main": {
			"main/site-1": graph.Edge[string]{
				Source: "main",
				Target: "main/site-1",
				Properties: graph.EdgeProperties{
					Attributes: map[string]string{},
					Weight:     0,
					Data:       interface{}(nil),
				},
			},
		},
		"main/site-1": {
			"main/site-1/site-component-1": graph.Edge[string]{
				Source: "main/site-1",
				Target: "main/site-1/site-component-1",
				Properties: graph.EdgeProperties{
					Attributes: map[string]string{},
					Weight:     0,
					Data:       interface{}(nil),
				},
			},
			"main/site-1/site-component-2": graph.Edge[string]{
				Source: "main/site-1",
				Target: "main/site-1/site-component-2",
				Properties: graph.EdgeProperties{
					Attributes: map[string]string{},
					Weight:     0,
					Data:       interface{}(nil),
				},
			},
		},
		"main/site-1/site-component-1": {},
		"main/site-1/site-component-2": {},
	}, am)
}

func TestToDependencyGraphInferred(t *testing.T) {
	cfg := &config.MachConfig{
		Filename: "main",
		MachComposer: config.MachComposer{
			Deployment: config.Deployment{
				Type: config.DeploymentSiteComponent,
			},
		},
		Sites: []config.SiteConfig{
			{
				Name: "site 1",
				Deployment: &config.Deployment{
					Type: config.DeploymentSiteComponent,
				},
				Identifier: "site-1",
				Components: []config.SiteComponentConfig{
					{
						Name: "site-component-1",
						Deployment: &config.Deployment{
							Type: config.DeploymentSiteComponent,
						},
					},
					{
						Name: "site-component-2",
						Deployment: &config.Deployment{
							Type: config.DeploymentSiteComponent,
						},
						Variables: variable.VariablesMap{
							"value-1": variable.MustCreateNewScalarVariable("${component.site-component-1.value-1}"),
						},
					},
				},
			},
		},
	}

	g, err := ToDependencyGraph(cfg, "")
	assert.NoError(t, err)

	am, _ := g.AdjacencyMap()
	assert.Equal(t, map[string]map[string]graph.Edge[string]{
		"main": {
			"main/site-1": graph.Edge[string]{
				Source:     "main",
				Target:     "main/site-1",
				Properties: graph.EdgeProperties{Attributes: map[string]string{}, Weight: 0, Data: interface{}(nil)},
			},
		},
		"main/site-1": {
			"main/site-1/site-component-1": graph.Edge[string]{
				Source:     "main/site-1",
				Target:     "main/site-1/site-component-1",
				Properties: graph.EdgeProperties{Attributes: map[string]string{}, Weight: 0, Data: interface{}(nil)},
			},
		},
		"main/site-1/site-component-1": {
			"main/site-1/site-component-2": graph.Edge[string]{
				Source:     "main/site-1/site-component-1",
				Target:     "main/site-1/site-component-2",
				Properties: graph.EdgeProperties{Attributes: map[string]string{}, Weight: 0, Data: interface{}(nil)},
			},
		},
		"main/site-1/site-component-2": {},
	}, am)
}

func TestToDependencyGraphCycleErr(t *testing.T) {
	c1 := &config.SiteComponentConfig{

		Name: "site-component-1",
		Deployment: &config.Deployment{
			Type: config.DeploymentSiteComponent,
		},
		DependsOnKeys: []string{
			"site-component-2",
		},
	}

	c2 := &config.SiteComponentConfig{
		Name: "site-component-2",
		Deployment: &config.Deployment{
			Type: config.DeploymentSiteComponent,
		},
		DependsOnKeys: []string{
			"site-component-1",
		},
	}

	c1.DependsOn = append(c1.DependsOn, c2)
	c2.DependsOn = append(c2.DependsOn, c1)

	cfg := &config.MachConfig{
		Filename: "main",
		MachComposer: config.MachComposer{
			Deployment: config.Deployment{
				Type: config.DeploymentSiteComponent,
			},
		},
		Sites: []config.SiteConfig{
			{
				Name:       "site 1",
				Identifier: "site-1",
				Deployment: &config.Deployment{
					Type: config.DeploymentSiteComponent,
				},
				Components: []config.SiteComponentConfig{*c1, *c2},
			},
		},
	}

	_, err := ToDependencyGraph(cfg, "")
	assert.Error(t, err)
	assert.IsType(t, &ValidationError{}, err)
	assert.Len(t, err.(*ValidationError).Errors, 1)
}
