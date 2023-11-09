package main

import (
	"bytes"
	"github.com/dominikbraun/graph/draw"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/goccy/go-graphviz"
)

var graphFlags struct {
	output           string
	outputDeployment string
}

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Print the execution graph for this project",
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
	graphCmd.Flags().StringVarP(&graphFlags.outputDeployment, "output-deployment", "", "./deployment-graph.png",
		"output file for the deployment image")
}

// TODO: turn both into single graph, where nested site components are also visible within the site terraform as a whole.
//
//	This should more clearly show the dependencies between the different sites and components.
func graphFunc(cmd *cobra.Command, _ []string) error {
	cfg := loadConfig(cmd, true)
	defer cfg.Close()

	g, err := dependency.FromConfig(cfg)
	if err != nil {
		return err
	}

	var buff bytes.Buffer
	err = draw.DOT(g.NodeGraph, &buff)
	if err != nil {
		return err
	}

	gv := graphviz.New()
	ggv, _ := graphviz.ParseBytes(buff.Bytes())

	if err := gv.RenderFilename(ggv, graphviz.PNG, graphFlags.output); err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msgf("Graph written to %s", graphFlags.output)

	dg, err := dependency.ToDeploymentGraph(g)
	if err != nil {
		return err
	}

	var buff2 bytes.Buffer
	err = draw.DOT(dg.NodeGraph, &buff2)
	if err != nil {
		return err
	}

	gv2 := graphviz.New()
	ggv2, _ := graphviz.ParseBytes(buff2.Bytes())

	if err := gv2.RenderFilename(ggv2, graphviz.PNG, graphFlags.outputDeployment); err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msgf("Deployment graph written to %s", graphFlags.outputDeployment)

	return nil
}
