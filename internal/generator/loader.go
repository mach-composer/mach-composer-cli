package generator

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"path"
)

type EmbedLoader struct {
	Content fs.FS
}

func (l *EmbedLoader) Abs(base, name string) string {
	return path.Join(path.Dir(base), name)
}

func (l *EmbedLoader) Get(p string) (io.Reader, error) {
	fullPath := path.Join("templates", p)
	f, err := l.Content.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(fmt.Errorf("error while closing %s: %v", fullPath, err))
		}
	}()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
