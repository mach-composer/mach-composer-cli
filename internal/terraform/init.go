package terraform

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func Init(ctx context.Context, path string) (string, error) {
	args := []string{"init"}

	return utils.RunTerraform(ctx, path, true, args...)
}
