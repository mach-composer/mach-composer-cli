package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type RunTerraformFunc func(ctx context.Context, catchOutput bool, cwd string, args ...string) (string, error)

var defaultRunTerraform RunTerraformFunc = utils.RunTerraform
