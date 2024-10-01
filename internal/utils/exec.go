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
	logger := log.With().
		Str("command", command).
		Strs("args", args).
		Str("cwd", cwd).Logger()

	logger.Debug().Msgf("Running: %s", command)

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Dir = cwd
	cmd.Env = os.Environ()

	//Currently keep the buffer in memory. We might want to change this to a file if the output is too large
	stdOut := new(bytes.Buffer)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
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
			//TODO: should we return the buffer here also?
			return "", fmt.Errorf("command (%s) failed: %w (args: %s , cwd: %s)", command, err, strings.Join(args, " "), cwd)
		}
	}

	logger.Debug().Msgf("Finished running: %s", command)

	return stdOut.String(), nil
}

func StopProcess(cmd *exec.Cmd) error {
	log.Info().Msg("Context cancelled, waiting for the command to exit...")
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
			log.Info().Msg("Waiting...")
			log.Info().Msg(cmd.ProcessState.String())
			if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
				log.Info().Msg("Command exited.")
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
