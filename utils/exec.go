package utils

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

func RunInteractive(ctx context.Context, command string, cwd string, args ...string) {
	logrus.Debugf("Running: %s %s\n", command, strings.Join(args, " "))

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
		log.Fatalln(err)
	}

	err = cmd.Wait()
	if err != nil {
		logrus.Fatalf("command exited: %s %s (in %s)", command, strings.Join(args, " "), cwd)
	}
}
