package config

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

// DecryptYaml takes a filename and returns the decrypted yaml.
// This command directly calls the sops binary instead of using the
// go.mozilla.org/sops/v3/decrypt package since that adds numerous dependencies
// and adds ~19mb to the generated binary
func DecryptYaml(filename string) ([]byte, error) {
	ctx := context.Background()

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return RunSops(ctx, wd, "-d", filename, "--output-type=yaml")
}

func RunSops(ctx context.Context, cwd string, args ...string) ([]byte, error) {
	logrus.Debugf("Running: sops %s\n", strings.Join(args, " "))
	cmd := exec.CommandContext(
		ctx,
		"sops",
		args...,
	)
	cmd.Dir = cwd
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Foreground: true,
	}
	return cmd.Output()
}
