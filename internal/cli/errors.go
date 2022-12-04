package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func PrintExitError(summary string, detail ...any) {
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	white := color.New(color.FgWhite, color.Bold).SprintFunc()
	fmt.Fprintln(os.Stderr, red("|"))
	fmt.Fprintln(os.Stderr, red("| Error:"), white(summary))
	fmt.Fprintln(os.Stderr, red("|"))
	if len(detail) > 0 {
		lines := []string{}
		for _, d := range detail {
			line := strings.TrimSpace(fmt.Sprintf("%s", d))
			parts := strings.Split(line, "\n")
			lines = append(lines, parts...)
		}

		for _, line := range lines {
			fmt.Fprintln(os.Stderr, red("|"), white(line))
		}
		fmt.Fprintln(os.Stderr, red("|"))
	}
	os.Exit(1)
}
