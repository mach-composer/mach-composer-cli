package terraform

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func Plan(ctx context.Context, path string, lock bool) (string, error) {
	cmd := []string{"plan"}

	if lock == false {
		cmd = append(cmd, "-lock=false")
	}

	cmd = append(cmd, fmt.Sprintf("-out=%s", PlanFile))
	return utils.RunTerraform(ctx, path, false, cmd...)
}
