package runner

import (
	"context"

	"github.com/labd/mach-composer/utils"
)

func RunTerraform(ctx context.Context, cwd string, args ...string) {
	utils.RunInteractive(ctx, "terraform", cwd, args...)
}
