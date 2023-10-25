package dependency

import (
	"fmt"
	"github.com/dominikbraun/graph"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"golang.org/x/exp/maps"
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

func (g *Graph) Vertices() Vertices {
	var vertices Vertices

	m, _ := g.AdjacencyMap()

	keys := maps.Keys(m)

	for _, k := range keys {
		v, _ := g.Vertex(k)
		vertices = append(vertices, v)
	}

	return vertices
}

// Routes determines all the possible paths between two nodes
func (g *Graph) Routes(source, target string) ([]Path, error) {
	var routes []Path

	m, err := g.PredecessorMap()
	if err != nil {
		return routes, err
	}

	eg := m[source]

	for _, pathElement := range eg {
		p := []string{pathElement.Source}
		newRoutes := fetchPathsToTarget(pathElement.Source, target, m, p)
		routes = append(routes, newRoutes...)
	}

	return routes, nil
}

func FromConfig(cfg *config.MachConfig) (*Graph, error) {
	var e = edges{}
	g := graph.New(func(n Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	project := &node{
		path:       strings.TrimSuffix(cfg.Filename, filepath.Ext(cfg.Filename)),
		identifier: cfg.Filename,
		typ:        ProjectType,
	}

	err := g.AddVertex(project)
	if err != nil {
		return nil, err
	}

	for _, siteConfig := range cfg.Sites {
		site := &node{
			path:       path.Join(project.Path(), siteConfig.Identifier),
			identifier: siteConfig.Identifier,
			typ:        SiteType,
			parent:     project,
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
			component := &node{
				path:       path.Join(site.Path(), componentConfig.Name),
				identifier: componentConfig.Name,
				typ:        SiteComponentType,
				parent:     site,
			}

			err = g.AddVertex(component)
			if err != nil {
				return nil, err
			}

			// First parse the explicit references. These always take precedence
			if dp := componentConfig.DependsOn; len(dp) > 0 {
				for _, dependency := range componentConfig.DependsOn {
					e.Add(component.Path(), path.Join(site.Path(), dependency))
				}
				continue
			}

			// If there are no explicit references, we need to check if there are any implicit ones
			if cp := componentConfig.Variables.ListComponents(); len(cp) > 0 {
				for _, dependency := range cp {
					e.Add(component.Path(), path.Join(site.Path(), dependency))
				}
			}
			if cp := componentConfig.Secrets.ListComponents(); len(cp) > 0 {
				for _, dependency := range cp {
					e.Add(component.Path(), path.Join(site.Path(), dependency))
				}
				continue
			}

			// Otherwise add the default link to the parent site
			e.Add(component.Path(), site.Path())
		}
	}

	// Process e
	var errList errorList
	for t, s := range e {
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
