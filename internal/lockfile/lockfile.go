package lockfile

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

var ErrLockFileNotFound = errors.New("Lockfile not found")

type LockFile struct {
	isNew bool `yaml:"-"`

	path string

	// Hash of all terraform files, alphabetically ordered
	TerraformHash string `yaml:"tf-hash"`

	// Hash of the mach-composer config file
	ConfigHash string `yaml:"cfg-hash"`

	Version string `yaml:"version"`
}

func (lf *LockFile) SetTerraformFiles() error {
	h, err := getTerraformFilesHash(lf.path)
	if err != nil {
		return err
	}
	lf.TerraformHash = h
	return nil
}

func (lf *LockFile) HasChanges(cfg *config.MachConfig) bool {
	if lf.isNew {
		return true
	}

	if cfg.ConfigHash != lf.ConfigHash {
		return true
	}

	md := cli.GetVersionMetadata()
	return md.Version != lf.Version
}

func (lf *LockFile) Write() error {
	data, err := yaml.Marshal(lf)
	if err != nil {
		return err
	}

	filename := filepath.Join(lf.path, ".mach-composer.lock")
	if err := utils.AFS.WriteFile(filename, data, 0600); err != nil {
		return err
	}

	return nil
}

func GetLockFile(cfg *config.MachConfig, path string) (*LockFile, error) {
	lockfile, err := readLockFile(path)
	if err != nil {
		if errors.Is(err, ErrLockFileNotFound) {
			metadata := cli.GetVersionMetadata()
			lockfile = &LockFile{
				path:       path,
				isNew:      true,
				ConfigHash: cfg.ConfigHash,
				Version:    metadata.Version,
			}
		}
	}

	return lockfile, nil
}

func readLockFile(path string) (*LockFile, error) {
	filename := filepath.Join(path, ".mach-composer.lock")
	data, err := utils.AFS.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrLockFileNotFound
		}
		return nil, err
	}

	ld := &LockFile{
		path: path,
	}
	if err := yaml.Unmarshal(data, ld); err != nil {
		return nil, err
	}

	return ld, nil
}

func getTerraformFilesHash(path string) (string, error) {
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
