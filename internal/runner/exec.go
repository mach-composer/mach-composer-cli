package runner

import (
	"context"
	"fmt"
	"os"

	"github.com/labd/mach-composer/internal/utils"
)

func RunTerraform(ctx context.Context, cwd string, args ...string) error {
	if _, err := os.Stat(cwd); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("The generated files are not found: %w", err)
		}
	}
	return utils.RunInteractive(ctx, "terraform", cwd, args...)
}
