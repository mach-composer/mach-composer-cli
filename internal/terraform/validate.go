package terraform

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func Validate(ctx context.Context, path string) (string, error) {
	cmd := []string{"validate"}

	return utils.RunTerraform(ctx, path, false, cmd...)
}
