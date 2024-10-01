package terraform

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"strings"
)

type ApplyOption func([]string) []string

func ApplyWithDestroy() ApplyOption {
	return func(args []string) []string {
		return append(args, "-destroy")
	}
}

func ApplyWithAutoApprove() ApplyOption {
	return func(args []string) []string {
		return append(args, "-auto-approve")
	}
}

func ApplyWithJson() ApplyOption {
	return func(args []string) []string {
		return append(args, "-json")
	}
}

func Apply(ctx context.Context, path string, opts ...ApplyOption) (string, error) {
	args := []string{"apply"}

	for _, opt := range opts {
		args = opt(args)
	}

	// If there is a plan then we should use it.
	planFilename, err := hasTerraformPlan(path)
	if err != nil {
		return "", err
	}
	if planFilename != "" {
		args = append(args, strings.TrimPrefix(planFilename, path+"/"))
	}

	return utils.RunTerraform(ctx, path, args...)
}
