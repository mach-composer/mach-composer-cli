package terraform

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func Init(ctx context.Context, path string, withBackend bool) (string, error) {
	args := []string{"init"}

	if !withBackend {
		args = append(args, "-backend=false")
	}

	return utils.RunTerraform(ctx, path, false, args...)
}
