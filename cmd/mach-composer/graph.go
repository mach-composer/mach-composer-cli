package main

import (
	"bytes"
	"fmt"
	"github.com/dominikbraun/graph/draw"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/spf13/cobra"
)

var graphFlags struct {
	output     string
	deployment bool
}

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Print the execution graph for this project",
	Long: `
Print the execution graph for this project. Note that the output will be in the DOT Language (https://graphviz.org/about/).

This output can be used to generate actual files with the dependency graph. 
A tool like graphviz can be used to make this transformation:
  
  'mach-composer graph -f main.yml | dot -Tpng -o image.png'
	
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return graphFunc(cmd, args)
	},
}

func init() {
	registerGenerateFlags(graphCmd)
	graphCmd.Flags().StringVarP(&graphFlags.output, "output", "", "./graph.png", "output file for the deployment image")
	graphCmd.Flags().BoolVarP(&graphFlags.deployment, "deployment", "d", false,
		"print the deployment graph instead of the dependency graph")
}

// TODO: turn both into single graph, where nested site components are also visible within the site terraform as a whole.
//
//	This should more clearly show the dependencies between the different sites and components.
func graphFunc(cmd *cobra.Command, _ []string) error {
	cfg := loadConfig(cmd, true)
	defer cfg.Close()

	g, err := dependency.ToDependencyGraph(cfg)
	if err != nil {
		return err
	}

	if graphFlags.deployment {
		dg, err := dependency.ToDeploymentGraph(cfg)
		if err != nil {
			return err
		}
		g = dg
	}

	var buff bytes.Buffer
	err = draw.DOT(g.NodeGraph, &buff)
	if err != nil {
		return err
	}

	fmt.Println(buff.String())

	return nil
}
