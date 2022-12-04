package utils

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

func RunSops(ctx context.Context, cwd string, args ...string) ([]byte, error) {
	log.Debug().Msgf("Running: sops %s\n", strings.Join(args, " "))
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
