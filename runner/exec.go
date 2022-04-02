package runner

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"syscall"

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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Foreground: true,
	}

	err := cmd.Run()

	if err != nil {
		logrus.Fatalf("terraform command exited: terraform %s (in %s)", strings.Join(args, " "), cwd)
	}
}
