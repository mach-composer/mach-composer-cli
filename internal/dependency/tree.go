package dependency

import (
	"fmt"
	"github.com/dominikbraun/graph"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"path"
	"slices"
)

type edges map[string][]string

func (e *edges) Add(to, from string) {
	if slices.Contains((*e)[to], from) {
		return
	}
	(*e)[to] = append((*e)[to], from)
}

func FromConfig(cfg *config.MachConfig) (graph.Graph[string, Node], Node, error) {
	var edges = edges{}
	g := graph.New(func(n Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	project := &Project{
		nodeImpl: nodeImpl{
			graph: g,
			path:  cfg.Filename,
		},
		Config: cfg,
	}

	err := g.AddVertex(project)
	if err != nil {
		return nil, nil, err
	}

	for _, siteConfig := range cfg.Sites {
		site := &Site{
			nodeImpl: nodeImpl{
				graph: g,
				path:  path.Join(project.Path(), siteConfig.Identifier),
			},
			Config: &siteConfig,
		}

		err = g.AddVertex(site)
		if err != nil {
			return nil, nil, err
		}

		err = g.AddEdge(project.Path(), site.Path())
		if err != nil {
			return nil, nil, err
		}

		for _, componentConfig := range siteConfig.Components {
			var cc = componentConfig
			component := &SiteComponent{
				nodeImpl: nodeImpl{
					graph: g,
					path:  path.Join(site.Path(), componentConfig.Name),
				},
				Config: &cc,
			}

			err = g.AddVertex(component)
			if err != nil {
				return nil, nil, err
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
		return nil, nil, &ValidationError{
			Msg:    "validation failed",
			Errors: errList,
		}
	}

	g, err = graph.TransitiveReduction(g)
	if err != nil {
		return nil, nil, err
	}

	return g, project, nil
}
