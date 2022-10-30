package runner

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetHash(path string) string {
	// Loop through all *.tf files
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	h := sha256.New()

	for _, f := range files {
		filename := filepath.Join(path, f.Name())
		if !strings.HasSuffix(filename, ".tf") {
			continue
		}
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(h, f); err != nil {
			_ = f.Close() // TODO: handle error
			log.Fatal(err)
		}

		_ = f.Close() // TODO: handle error
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
