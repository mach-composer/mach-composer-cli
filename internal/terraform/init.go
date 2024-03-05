package terraform

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func Init(ctx context.Context, path string) (string, error) {
	if !terraformIsInitialized(path) {
		return utils.RunTerraform(ctx, path, true, "init")
	}
	return "", nil
}
