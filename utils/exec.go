package utils

import (
	"context"
	"io"
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

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalln(err)
	}

	go io.Copy(os.Stderr, stderr)
	go io.Copy(os.Stdout, stdout)
	go io.Copy(stdin, os.Stdin)

	err = cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}

	err = cmd.Wait()
	if err != nil {
		logrus.Fatalf("command exited: %s %s (in %s)", command, strings.Join(args, " "), cwd)
	}
}
