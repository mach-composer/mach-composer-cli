package cli

import (
	"context"
	"io"
	"os"
)

const WriterKey = "log-writer"

func ContextWithLogWriter(ctx context.Context, w io.Writer) context.Context {
	return context.WithValue(ctx, WriterKey, w)
}

// LogWriterFromContext returns the currently configured log writer for reuse in routine loggers
func LogWriterFromContext(ctx context.Context) io.Writer {
	if v := ctx.Value(WriterKey); v != nil {
		return v.(io.Writer)
	}

	return os.Stdout
}

const OutputKey = "output"

type OutputType string

const (
	OutputTypeConsole OutputType = "console"
	OutputTypeJSON    OutputType = "json"
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
