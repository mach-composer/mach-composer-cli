package utils

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

func RunSops(ctx context.Context, cwd string, args ...string) ([]byte, error) {
	logrus.Debugf("Running: sops %s\n", strings.Join(args, " "))
	cmd := exec.CommandContext(
		ctx,
		"sops",
		args...,
	)
	cmd.Dir = cwd
	cmd.Stderr = os.Stderr
	CmdSetForeground(cmd)
	return cmd.Output()
}
