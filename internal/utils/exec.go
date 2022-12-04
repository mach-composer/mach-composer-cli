package utils

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

func RunInteractive(ctx context.Context, command string, cwd string, args ...string) error {
	log.Debug().Msgf("Running: %s %s\n", command, strings.Join(args, " "))

	cmd := exec.CommandContext(
		ctx,
		command,
		args...,
	)
	cmd.Dir = cwd
	cmd.Env = os.Environ()

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("command exited: %s %s (in %s)", command, strings.Join(args, " "), cwd)
	}
	return nil
}
