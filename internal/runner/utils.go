package runner

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

func terraformIsInitialized(ctx context.Context, path string) bool {
	tfLockFile := filepath.Join(path, ".terraform.lock.hcl")
	if _, err := os.Stat(tfLockFile); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Ctx(ctx).Fatal().Err(err)
	}
	return true
}

func terraformCanPlan(ctx context.Context, n graph.Node) (bool, error) {
	parents, err := n.Parents()
	if err != nil {
		return false, err
	}

	// Sites can always plan, so no need to check
	if n.Type() == graph.SiteType {
		return true, nil
	}

	for _, parent := range parents {
		if parent.Type() == graph.SiteType {
			sp, ok := parent.(*graph.Site)
			if !ok {
				return false, fmt.Errorf("failed to cast parent to site")
			}

			if len(sp.NestedNodes) == 0 {
				log.Debug().Msgf("site node does not contain components, so no output is available")
				continue
			}
		}

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
