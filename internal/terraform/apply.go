package terraform

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"strings"
)

func Apply(ctx context.Context, path string, destroy, autoApprove bool) (string, error) {
	args := []string{"apply"}

	if destroy {
		args = append(args, "-destroy")
	}

	if autoApprove {
		args = append(args, "-auto-approve")
	}

	if ctx.Value(cli.OutputKey) == cli.OutputTypeJSON {
		args = append(args, "-json")
	}

	// If there is a plan then we should use it.
	planFilename, err := hasTerraformPlan(path)
	if err != nil {
		return "", err
	}
	if planFilename != "" {
		args = append(args, strings.TrimPrefix(planFilename, path+"/"))
	}

	return utils.RunTerraform(ctx, path, true, args...)
}
