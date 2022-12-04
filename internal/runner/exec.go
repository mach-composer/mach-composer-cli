package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/labd/mach-composer/internal/utils"
)

func RunTerraform(ctx context.Context, cwd string, args ...string) error {
	if _, err := os.Stat(cwd); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("The generated files are not found: %w", err)
		}
	}

	execPath, err := exec.LookPath("terraform")
	if err != nil {
		return err
	}

	return utils.RunInteractive(ctx, execPath, cwd, args...)
}
