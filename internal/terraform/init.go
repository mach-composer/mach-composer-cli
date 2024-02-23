package terraform

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func Init(ctx context.Context, path string) (string, error) {
	if !terraformIsInitialized(path) {
		return utils.RunTerraform(ctx, path, true, "init")
	}
	return fmt.Sprintf("initialization skipped on path %s", path), nil
}
