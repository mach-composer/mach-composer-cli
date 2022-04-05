package runner

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/labd/mach-composer/utils"
	"github.com/sirupsen/logrus"
)

func RunTerraform(ctx context.Context, cwd string, args ...string) {
	logrus.Debugf("Running: terraform %s\n", strings.Join(args, " "))
	cmd := exec.CommandContext(
		ctx,
		"terraform",
		args...,
	)
	cmd.Dir = cwd

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = os.Environ()
	utils.CmdSetForegrond(cmd)

	err := cmd.Run()
	if err != nil {
		logrus.Fatalf("terraform command exited: terraform %s (in %s)", strings.Join(args, " "), cwd)
	}
}
