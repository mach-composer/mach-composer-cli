package lockfile

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mach-composer/mach-composer-cli/internal/cli"
)

const FileName = ".mach-composer.lock"

type LockFile struct {
	isNew bool   `yaml:"-"`
	path  string `yaml:"-"`

	// Hash of all terraform files, alphabetically ordered
	TerraformHash string `yaml:"tf-hash"`

	// Hash of the mach-composer config file
	ConfigHash string `yaml:"cfg-hash"`

	Version string `yaml:"version"`
}

func (lf *LockFile) HasChanges(hash string) bool {
	if lf.isNew {
		return true
	}

	if hash != lf.ConfigHash {
		return true
	}

	md := cli.GetVersionMetadata()
	return md.Version != lf.Version
}

func (lf *LockFile) Update(hash string) error {
	lf.ConfigHash = hash
	lf.Version = cli.GetVersionMetadata().Version
	tfHash, err := getTerraformFilesHash(lf.path)
	if err != nil {
		return err
	}
	lf.TerraformHash = tfHash

	return nil
}

func getTerraformFilesHash(path string) (string, error) {
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
