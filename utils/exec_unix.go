//go:build !windows

package utils

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/mattn/go-isatty"
)

func CmdSetForeground(cmd *exec.Cmd) {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Foreground: true,
		}
	}
}
