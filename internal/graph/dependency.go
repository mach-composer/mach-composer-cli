package graph

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

type edgeSets map[string][]string

func (e *edgeSets) Add(to, from string) {
	if slices.Contains((*e)[to], from) {
		return
	}
	(*e)[to] = append((*e)[to], from)
}

// ToDependencyGraph will transform a MachConfig into a graph of dependencies connected by different relations
func ToDependencyGraph(cfg *config.MachConfig, outPath string) (*Graph, error) {
	var edges = edgeSets{}
	g := graph.New(func(n Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	projectIdentifier := strings.TrimSuffix(cfg.Filename, filepath.Ext(cfg.Filename))

	p := path.Join(outPath, projectIdentifier)

	project := NewProject(g, p, projectIdentifier, cfg.MachComposer.Deployment.Type, *cfg)

	err := g.AddVertex(project)
	if err != nil {
		return nil, err
	}

	for _, siteConfig := range cfg.Sites {
		p = path.Join(project.Path(), siteConfig.Identifier)
		site := NewSite(g, p,
			siteConfig.Identifier,
			siteConfig.Deployment.Type,
			project,
			*cfg,
			siteConfig,
		)

		err = g.AddVertex(site)
		if err != nil {
			return nil, err
		}

		err = g.AddEdge(project.Path(), site.Path())
		if err != nil {
			return nil, err
		}

		for _, componentConfig := range siteConfig.Components {
			log.Debug().Msgf("Deploying site component %s separately",
				componentConfig.Name)

			p = path.Join(site.Path(), componentConfig.Name)
			component := NewSiteComponent(g, p,
				componentConfig.Name,
				componentConfig.Deployment.Type,
				site,
				*cfg,
				siteConfig,
				componentConfig,
			)

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
			if cp := componentConfig.Variables.ListReferencedComponents(); len(cp) > 0 {
				for _, dependency := range cp {
					edges.Add(component.Path(), path.Join(site.Path(), dependency))
				}
			}
			if cp := componentConfig.Secrets.ListReferencedComponents(); len(cp) > 0 {
				for _, dependency := range cp {
					edges.Add(component.Path(), path.Join(site.Path(), dependency))
				}
				continue
			}

			// Otherwise add the default link to the ancestor site
			edges.Add(component.Path(), site.Path())
		}
	}

	// Process edges
	var errList errorList
	for target, sources := range edges {
		for _, source := range sources {
			err = g.AddEdge(source, target)
			if err != nil {
				errList.AddError(fmt.Errorf("failed to add dependency from %v to %v: %w", source, target, err))
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

	return &Graph{Graph: g, StartNode: project}, nil
}
