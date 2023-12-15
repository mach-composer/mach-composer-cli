package terraform

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func Show(ctx context.Context, path string, noColor bool) (string, error) {
	filename, err := hasTerraformPlan(path)
	if err != nil {
		return "", err
	}
	if filename == "" {
		return "", fmt.Errorf("no plan found for path %s. Did you run `mach-composer plan`", path)
	}

	cmd := []string{"show", filename}
	if noColor {
		cmd = append(cmd, "-no-color")
	}
	return utils.RunTerraform(ctx, path, false, cmd...)
}
