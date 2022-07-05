package runner

import (
	"context"
	"fmt"
	"os"

	"github.com/labd/mach-composer/internal/utils"
)

func RunTerraform(ctx context.Context, cwd string, args ...string) {
	if _, err := os.Stat(cwd); err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "The generated files are not found. Did you run mach-composer generate?")
			os.Exit(1)
		}
	}

	utils.RunInteractive(ctx, "terraform", cwd, args...)
}
