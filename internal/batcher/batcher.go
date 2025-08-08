package batcher

import (
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"slices"
)

type BatchFunc func(g *graph.Graph) (map[int][]graph.Node, error)

type Batcher string

func Factory(cfg *config.MachConfig) (BatchFunc, error) {
	switch cfg.MachComposer.Batcher.Type {
	case "":
		fallthrough
	case "simple":
		return simpleBatchFunc(), nil
	case "site":
		var siteOrder, err = DetermineSiteOrder(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed determining site order: %w", err)
		}

		return siteBatchFunc(siteOrder), nil
	default:
		return nil, fmt.Errorf("unknown batch type %s", cfg.MachComposer.Batcher.Type)
	}
}

func DetermineSiteOrder(cfg *config.MachConfig) ([]string, error) {
	var identifiers = cfg.Sites.Identifiers()
	var siteOrder = make([]string, len(identifiers))

	if len(cfg.MachComposer.Batcher.SiteOrder) > 0 {
		// Use the site order from the configuration if provided
		siteOrder = cfg.MachComposer.Batcher.SiteOrder

		// Make sure the site order contains the same fields as the identifiers
		if len(siteOrder) != len(identifiers) {
			return nil, fmt.Errorf("site order length %d does not match identifiers length %d", len(siteOrder), len(identifiers))
		}
		for _, siteIdentifier := range siteOrder {
			if !slices.Contains(identifiers, siteIdentifier) {
				return nil, fmt.Errorf("site order contains siteIdentifier %s that is not in the identifiers list", siteIdentifier)
			}
		}

	} else {
		for i, identifier := range identifiers {
			siteOrder[i] = identifier
		}
	}

	return siteOrder, nil
}
