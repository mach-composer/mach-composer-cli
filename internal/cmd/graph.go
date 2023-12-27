package cmd

import (
	"bytes"
	"fmt"
	"github.com/dominikbraun/graph/draw"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
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
		preprocessCommonFlags(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return graphFunc(cmd, args)
	},
}

func init() {
	registerCommonFlags(graphCmd)
	graphCmd.Flags().StringVarP(&graphFlags.output, "output", "", "./graph.png", "output file for the deployment image")
	graphCmd.Flags().BoolVarP(&graphFlags.deployment, "deployment", "d", false,
		"print the deployment graph instead of the dependency graph")
}

func graphFunc(cmd *cobra.Command, _ []string) error {
	cfg := loadConfig(cmd, true)
	defer cfg.Close()

	g, err := graph.ToDependencyGraph(cfg, commonFlags.outputPath)
	if err != nil {
		return err
	}

	if graphFlags.deployment {
		dg, err := graph.ToDeploymentGraph(cfg, commonFlags.outputPath)
		if err != nil {
			return err
		}
		g = dg
	}

	var buff bytes.Buffer
	err = draw.DOT(g.Graph, &buff)
	if err != nil {
		return err
	}

	fmt.Println(buff.String())

	return nil
}
