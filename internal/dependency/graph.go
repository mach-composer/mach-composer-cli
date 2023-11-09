package dependency

import (
	"fmt"
	"github.com/dominikbraun/graph"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
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

type NodeGraph = graph.Graph[string, Node]

type Graph struct {
	NodeGraph
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

func ToDependencyGraph(cfg *config.MachConfig) (*Graph, error) {
	var edges = edgeSets{}
	g := graph.New(func(n Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	project := &Project{
		node: node{
			path:           strings.TrimSuffix(cfg.Filename, filepath.Ext(cfg.Filename)),
			identifier:     cfg.Filename,
			typ:            ProjectType,
			deploymentType: cfg.MachComposer.Deployment.Type,
		},
		ProjectConfig: cfg,
	}

	err := g.AddVertex(project)
	if err != nil {
		return nil, err
	}

	for _, siteConfig := range cfg.Sites {
		site := &Site{
			node: node{
				path:           path.Join(project.Path(), siteConfig.Identifier),
				identifier:     siteConfig.Identifier,
				typ:            SiteType,
				parent:         project,
				deploymentType: siteConfig.Deployment.Type,
			},
			SiteConfig: siteConfig,
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
			var p = path.Join(site.Path(), componentConfig.Name)

			log.Debug().Msgf("Deploying site component %s separately", componentConfig.Name)
			component := &SiteComponent{
				node: node{
					path:           p,
					identifier:     componentConfig.Name,
					typ:            SiteComponentType,
					parent:         site,
					deploymentType: componentConfig.Deployment.Type,
				},
				SiteConfig:          siteConfig,
				SiteComponentConfig: componentConfig,
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

			// Otherwise add the default link to the parent site
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

	return &Graph{NodeGraph: g, StartNode: project}, nil
}

func validateDeployment(g NodeGraph, start string) error {
	var errList errorList
	err := graph.DFS(g, start, func(p string) bool {
		n, _ := g.Vertex(p)

		am, _ := g.AdjacencyMap()
		edges := am[p]

		for _, edge := range edges {
			child, _ := g.Vertex(edge.Target)
			//TODO: this is a weird check, maybe we want to make it more agnostic of what type of node it is
			if n.Type() == SiteComponentType && n.Independent() && !child.Independent() {
				errList.AddError(fmt.Errorf("node %s is independent but has a dependent child %s", n.Path(), child.Path()))
			}
		}

		return false
	})
	if err != nil {
		return err
	}
	if len(errList) > 0 {
		return &ValidationError{
			Msg:    "validation failed",
			Errors: errList,
		}
	}

	return nil
}
