package runner

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func getFileHash(path string) (string, error) {
	// Loop through all *.tf files
	files, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	h := sha256.New()

	for _, f := range files {
		filename := filepath.Join(path, f.Name())
		if !strings.HasSuffix(filename, ".tf") {
			continue
		}
		f, err := os.Open(filename)
		if err != nil {
			return "", err
		}
		defer f.Close()

		if _, err := io.Copy(h, f); err != nil {
			return "", err
		}
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
