package terraform

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type PlanOption func([]string) []string

func PlanWithNoLock() PlanOption {
	return func(args []string) []string {
		return append(args, "-lock=false")
	}
}

func PlanWithJson() PlanOption {
	return func(args []string) []string {
		return append(args, "-json")
	}
}

func Plan(ctx context.Context, path string, opts ...PlanOption) (string, error) {
	args := []string{"plan"}

	for _, opt := range opts {
		args = opt(args)
	}

	args = append(args, fmt.Sprintf("-out=%s", PlanFile))

	return utils.RunTerraform(ctx, path, args...)
}
