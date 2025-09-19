package runner

import (
	"context"
	"slices"
	"strings"

	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/mach-composer/mach-composer-cli/internal/hash"
)

type Filter string

func (f Filter) WithDependencies() bool {
	return strings.HasSuffix(string(f), "...")
}

func (f Filter) Base() string {
	return strings.TrimSuffix(string(f), "...")
}

type Filters []Filter

func (fs Filters) Find(f string) *Filter {
	for _, filter := range fs {
		if filter.Base() == f {
			return &filter
		}
		el := strings.Split(f, "/")
		if slices.Contains(el, filter.Base()) {
			return &filter
		}
	}
	return nil
}

func (fs Filters) Contains(v string) bool {
	f := fs.Find(v)
	return f != nil
}

func toFilters(s []string) Filters {
	var filters Filters
	for _, f := range s {
		filters = append(filters, Filter(f))
	}
	return filters
}

func childrenAreActive(n graph.Node, f *Filter, previousParentActive bool) bool {
	return (n.Active() && f != nil && f.WithDependencies()) || previousParentActive
}

func filterNode(ctx context.Context, hashFetcher hash.Handler, g *graph.Graph, path string, filters Filters, parentActive bool) error {
	n, err := g.Vertex(path)
	if err != nil {
		return err
	}

	am, err := g.AdjacencyMap()
	if err != nil {
		return err
	}

	f := filters.Find(n.Identifier())
	if f != nil || parentActive {
		n.SetActive(true)
	}

	for _, child := range am[path] {
		if err = filterNode(ctx, hashFetcher, g, child.Target, filters, childrenAreActive(n, f, parentActive)); err != nil {
			return err
		}
	}

	return nil
}

func filterGraph(ctx context.Context, g *graph.Graph, hashFetcher hash.Handler, filters Filters) error {
	if len(filters) == 0 {
		for _, v := range g.Vertices() {
			v.SetActive(true)
		}
		return nil
	}

	return filterNode(ctx, hashFetcher, g, g.StartNode.Path(), filters, false)
}
