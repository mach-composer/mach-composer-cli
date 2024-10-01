package terraform

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type ShowOption func([]string) []string

func ShowWithNoColor() ShowOption {
	return func(args []string) []string {
		return append(args, "-no-color")
	}
}

func ShowWithJson() ShowOption {
	return func(args []string) []string {
		return append(args, "-json")
	}
}

func Show(ctx context.Context, path string, opts ...ShowOption) (string, error) {
	filename, err := hasTerraformPlan(path)
	if err != nil {
		return "", err
	}
	if filename == "" {
		return "", fmt.Errorf("no plan found for path %s. Did you run `mach-composer plan`", path)
	}

	args := []string{"show", filename}

	for _, opt := range opts {
		args = opt(args)
	}

	return utils.RunTerraform(ctx, path, args...)
}
