//go:build !windows

package utils

import (
	"os/exec"
	"syscall"
)

func CmdSetForegrond(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Foreground: true,
	}
}
