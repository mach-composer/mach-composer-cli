package main

import (
	"fmt"
	"github.com/dominikbraun/graph"
	"github.com/mach-composer/mach-composer-cli/internal/tree"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
	"strings"
)

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Determine the execution graph for the project",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return graphFunc(cmd, args)
	},
}

func init() {
	registerGenerateFlags(graphCmd)
}

func graphFunc(cmd *cobra.Command, _ []string) error {
	cfg := loadConfig(cmd, true)
	defer cfg.Close()

	g, root, err := tree.FromConfig(cfg)
	if err != nil {
		return err
	}

	am, err := g.AdjacencyMap()
	if err != nil {
		return err
	}

	tr := treeprint.NewWithRoot(root)
	ram, _ := am[root.Path()]

	for k, _ := range ram {
		processNode(tr, am, k)
	}

	fmt.Println(tr.String())

	return nil
}

func processNode(tp treeprint.Tree, am map[string]map[string]graph.Edge[string], path string) {
	ram := am[path]

	e := strings.Split(path, "/")

	if len(ram) == 0 {
		tp.AddNode(e[len(e)-1])
		return
	}

	tp2 := tp.AddBranch(e[len(e)-1])

	for k, _ := range ram {
		processNode(tp2, am, k)
	}
}
