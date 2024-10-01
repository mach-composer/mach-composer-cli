package terraform

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type ValidateOption func([]string) []string

func ValidateWithJson() ValidateOption {
	return func(args []string) []string {
		return append(args, "-json")
	}
}

func Validate(ctx context.Context, path string, opts ...ValidateOption) (string, error) {
	args := []string{"validate"}

	for _, opt := range opts {
		args = opt(args)
	}

	return utils.RunTerraform(ctx, path, args...)
}
