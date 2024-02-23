package terraform

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func Plan(ctx context.Context, path string, lock bool) (string, error) {
	args := []string{"plan"}

	if lock == false {
		args = append(args, "-lock=false")
	}

	if ctx.Value(cli.OutputKey) == cli.OutputTypeJSON {
		args = append(args, "-json")
	}

	args = append(args, fmt.Sprintf("-out=%s", PlanFile))
	return utils.RunTerraform(ctx, path, true, args...)
}
