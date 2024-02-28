package terraform

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func Init(ctx context.Context, path string) error {
	if !terraformIsInitialized(path) {
		if _, err := utils.RunTerraform(ctx, path, false, "init"); err != nil {
			return err
		}
	}
	return nil
}
