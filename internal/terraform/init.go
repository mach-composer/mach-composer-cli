package terraform

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type InitOption func([]string) []string

func InitWithDisableBackend() InitOption {
	return func(args []string) []string {
		return append(args, "--backend=false")
	}
}

func Init(ctx context.Context, path string, opts ...InitOption) (string, error) {
	args := []string{"init"}

	for _, opt := range opts {
		args = opt(args)
	}

	return utils.RunTerraform(ctx, path, args...)
}
