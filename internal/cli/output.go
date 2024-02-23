package cli

import (
	"context"
	"fmt"
)

const OutputKey = "output"

type OutputType string

const (
	OutputTypeConsole OutputType = "console"
	OutputTypeJSON    OutputType = "json"
	OutputTypeGitHub  OutputType = "github"
)

func ContextWithOutput(ctx context.Context, output OutputType) context.Context {
	return context.WithValue(ctx, OutputKey, output)
}

// OutputFromContext returns the output type from the context. If no output
// type is set console is returned
func OutputFromContext(ctx context.Context) OutputType {
	if v := ctx.Value(OutputKey); v != nil {
		return v.(OutputType)
	}
	return OutputTypeConsole
}

func ConvertOutputType(s string) (OutputType, error) {
	switch s {
	case "console":
		return OutputTypeConsole, nil
	case "json":
		return OutputTypeJSON, nil
	case "github":
		return OutputTypeGitHub, nil
	default:
		return "", fmt.Errorf("unknown output type: %s", s)
	}
}
