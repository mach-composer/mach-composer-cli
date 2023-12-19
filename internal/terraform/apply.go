package terraform

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"strings"
)

func Apply(ctx context.Context, path string, destroy, autoApprove bool) (string, error) {
	cmd := []string{"apply"}

	if destroy {
		cmd = append(cmd, "-destroy")
	}

	if autoApprove {
		cmd = append(cmd, "-auto-approve")
	}

	// If there is a plan then we should use it.
	planFilename, err := hasTerraformPlan(path)
	if err != nil {
		return "", err
	}
	if planFilename != "" {
		cmd = append(cmd, strings.TrimPrefix(planFilename, path+"/"))
	}

	return utils.RunTerraform(ctx, path, false, cmd...)
}
