package terraform

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/lockfile"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func Init(ctx context.Context, hash, path string) error {
	lf, err := lockfile.GetLock(hash, path)
	if err != nil {
		return err
	}

	if !terraformIsInitialized(path) || lf.HasChanges(hash) {
		if _, err = utils.RunTerraform(ctx, path, false, "init"); err != nil {
			return err
		}
	}
	return nil
}
