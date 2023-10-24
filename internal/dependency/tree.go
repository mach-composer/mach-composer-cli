package dependency

import (
	"fmt"
	"github.com/dominikbraun/graph"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/rs/zerolog/log"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

type edges map[string][]string

func (e *edges) Add(to, from string) {
	if slices.Contains((*e)[to], from) {
		return
	}
	(*e)[to] = append((*e)[to], from)
}

type Graph struct {
	graph.Graph[string, Node]
	StartNode Node
}

func (g *Graph) Walk(fn func(n Node) error) error {
	return graph.DFS[string, Node](g, g.StartNode.Path(), func(k string) bool {
		v, _ := g.Vertex(k)
		err := fn(v)
		if err != nil {
			log.Error().Msgf("Could not process node %s: %s", k, err.Error())
		}
		return false
	})
}

func FromConfig(cfg *config.MachConfig) (*Graph, error) {
	var edges = edges{}
	g := graph.New(func(n Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	project := &Project{
		nodeImpl: nodeImpl{
			graph:      g,
			path:       strings.TrimSuffix(cfg.Filename, filepath.Ext(cfg.Filename)),
			identifier: cfg.Filename,
		},
		Config: cfg,
	}

	err := g.AddVertex(project)
	if err != nil {
		return nil, err
	}

	for _, siteConfig := range cfg.Sites {
		site := &Site{
			nodeImpl: nodeImpl{
				graph:      g,
				path:       path.Join(project.Path(), siteConfig.Identifier),
				identifier: siteConfig.Identifier,
			},
			ProjectConfig: cfg,
			Config:        &siteConfig,
		}

		err = g.AddVertex(site)
		if err != nil {
			return nil, err
		}

		err = g.AddEdge(project.Path(), site.Path())
		if err != nil {
			return nil, err
		}

		for _, componentConfig := range siteConfig.Components {
			var cc = componentConfig
			component := &SiteComponent{
				nodeImpl: nodeImpl{
					graph:      g,
					path:       path.Join(site.Path(), componentConfig.Name),
					identifier: componentConfig.Name,
				},
				ProjectConfig: cfg,
				SiteConfig:    &siteConfig,
				Config:        &cc,
			}

			err = g.AddVertex(component)
			if err != nil {
				return nil, err
			}

			// First parse the explicit references. These always take precedence
			if dp := componentConfig.DependsOn; len(dp) > 0 {
				for _, dependency := range componentConfig.DependsOn {
					edges.Add(component.Path(), path.Join(site.Path(), dependency))
				}
				continue
			}

			// If there are no explicit references, we need to check if there are any implicit ones
			if cp := componentConfig.Variables.ListComponents(); len(cp) > 0 {
				for _, dependency := range cp {
					edges.Add(component.Path(), path.Join(site.Path(), dependency))
				}
			}
			if cp := componentConfig.Secrets.ListComponents(); len(cp) > 0 {
				for _, dependency := range cp {
					edges.Add(component.Path(), path.Join(site.Path(), dependency))
				}
				continue
			}

			// Otherwise add the default link to the parent site
			edges.Add(component.Path(), site.Path())
		}
	}

	// Process edges
	var errList errorList
	for t, s := range edges {
		for _, f := range s {
			err = g.AddEdge(f, t)
			if err != nil {
				errList.AddError(fmt.Errorf("failed to add dependency from %v to %v: %w", f, t, err))
			}
		}
	}

	if len(errList) > 0 {
		return nil, &ValidationError{
			Msg:    "validation failed",
			Errors: errList,
		}
	}

	g, err = graph.TransitiveReduction(g)
	if err != nil {
		return nil, err
	}

	return &Graph{
		Graph:     g,
		StartNode: project,
	}, nil
}
