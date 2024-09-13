package batcher

import "github.com/mach-composer/mach-composer-cli/internal/graph"

type BatchFunc func(g *graph.Graph) map[int][]graph.Node
