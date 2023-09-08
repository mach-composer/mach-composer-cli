package main

import (
	"bytes"
	"fmt"
	"github.com/dominikbraun/graph/draw"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/spf13/cobra"
)

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
}

func graphFunc(cmd *cobra.Command, _ []string) error {
	cfg := loadConfig(cmd, true)
	defer cfg.Close()

	t, _, err := dependency.FromConfig(cfg)
	if err != nil {
		return err
	}

	var buff bytes.Buffer
	err = draw.DOT(t, &buff, draw.GraphAttribute("label", cfg.Filename))
	if err != nil {
		return err
	}

	fmt.Println(buff.String())

	return nil
}
