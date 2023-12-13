package runner

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/rs/zerolog/log"
)

type PlanOptions struct {
	Lock    bool
	Workers int
}

func TerraformPlan(ctx context.Context, dg *dependency.Graph, opts *PlanOptions) error {
	if err := batchRun(ctx, dg, opts.Workers, func(ctx context.Context, n dependency.Node) (string, error) {
		return terraformPlan(ctx, n, n.Path(), opts.Lock)
	}); err != nil {
		return err
	}

	return nil
}

func terraformPlan(ctx context.Context, n dependency.Node, path string, lock bool) (string, error) {
	cmd := []string{"plan"}

	if n.Type() == dependency.SiteComponentType {
		parents, err := n.Parents()
		if err != nil {
			return "", err
		}

		var skip = false
		var details []string
		for _, p := range parents {
			if p.Outputs().Type().HasAttribute("output") == false {
				details = append(details, fmt.Sprintf("parent %s has no output attribute", p.Path()))
				skip = true
			}
		}
		if skip == true {
			log.Warn().Strs("state", details).Msgf("Skipping plan for %s because of missing parent states", n.Path())
			return "", nil
		}
	}

	if lock == false {
		cmd = append(cmd, "-lock=false")
	}

	cmd = append(cmd, fmt.Sprintf("-out=%s", PlanFile))
	return utils.RunTerraform(ctx, false, path, cmd...)
}
