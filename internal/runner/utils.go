package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

func terraformIsInitialized(path string) bool {
	tfLockFile := filepath.Join(path, ".terraform.lock.hcl")
	if _, err := os.Stat(tfLockFile); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatal().Err(err)
	}
	return true
}

func terraformCanPlan(ctx context.Context, n graph.Node) (bool, error) {
	parents, err := n.Parents()
	if err != nil {
		return false, err
	}

	for _, parent := range parents {
		v, err := utils.GetTerraformOutputs(ctx, parent.Path())
		if err != nil {
			return false, nil
		}
		a := v.Type().AttributeTypes()
		if len(a) == 0 {
			return false, nil
		}
	}
	return true, nil
}
