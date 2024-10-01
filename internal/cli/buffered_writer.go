package cli

import (
	"io"
)

// BufferedWriter is a writer that buffers all writes until Flush is called. This is useful in go routines where we
// don't want the logs to be mixed up between routines
type BufferedWriter struct {
	io.Writer
	buffer []string
}

func NewBufferedWriter(w io.Writer) *BufferedWriter {
	return &BufferedWriter{Writer: w}
}

func (b *BufferedWriter) Write(p []byte) (n int, err error) {
	b.buffer = append(b.buffer, string(p))
	return len(p), nil
}

// Flush writes all buffered data to the inner writer
func (b *BufferedWriter) Flush() error {
	for _, p := range b.buffer {
		if _, err := b.Writer.Write([]byte(p)); err != nil {
			return err
		}
	}
	return nil
}
