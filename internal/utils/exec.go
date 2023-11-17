package utils

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func RunInteractive(ctx context.Context, command string, cwd string, args ...string) (string, error) {
	log.Debug().Msgf("Running: %s %s\n", command, strings.Join(args, " "))

	cmd := exec.CommandContext(
		ctx,
		command,
		args...,
	)
	cmd.Dir = cwd
	cmd.Env = os.Environ()

	stdOut := new(bytes.Buffer)
	stdErr := new(bytes.Buffer)

	cmd.Stdin = os.Stdin
	cmd.Stderr = stdErr
	cmd.Stdout = stdOut

	err := cmd.Start()
	if err != nil {
		return "", err
	}

	// Wait for the command to complete or the context to be cancelled
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		return "", StopProcess(cmd)

	case err := <-done:
		if err != nil {
			return "", fmt.Errorf("command (%s) failed: %w (args: %s , cwd: %s): %s", command, err, strings.Join(args, " "),
				cwd, stdErr.String())
		}
	}

	return stdOut.String(), nil
}

func StopProcess(cmd *exec.Cmd) error {
	fmt.Println("Context cancelled, waiting for the command to exit...")
	err := cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return fmt.Errorf("failed to send interrupt signal: %w\n", err)
	}

	// Poll for the process to exit
	pollInterval := 250 * time.Millisecond
	ticker := time.NewTicker(pollInterval)
	current := time.Now()
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Waiting...")
			fmt.Println(cmd.ProcessState)
			if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
				fmt.Println("Command exited.")
				return nil
			}
			if current.Add(10 * time.Second).After(time.Now()) {
				err := cmd.Process.Signal(os.Kill)
				if err != nil {
					return fmt.Errorf("failed to send interrupt signal: %w\n", err)
				}
				return fmt.Errorf("command did not exit within the specified duration, killing")
			}
		}
	}
}
