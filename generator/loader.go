package generator

import (
	"bytes"
	"io"
	"io/fs"
	"io/ioutil"
	"path/filepath"
)

type EmbedLoader struct {
	Content fs.FS
}

func (l *EmbedLoader) Abs(base, name string) string {
	return filepath.Join(filepath.Dir(base), name)
}

func (l *EmbedLoader) Get(path string) (io.Reader, error) {
	fullPath := filepath.Join("templates", path)
	f, err := l.Content.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
