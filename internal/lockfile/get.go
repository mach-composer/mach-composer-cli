package lockfile

import (
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func GetLock(hash, path string) (*LockFile, error) {
	lockfile, exists, err := readLockFile(path, FileName)
	if err != nil {
		return nil, err
	}
	if !exists {
		metadata := cli.GetVersionMetadata()
		lockfile = &LockFile{
			path:       path,
			isNew:      true,
			ConfigHash: hash,
			Version:    metadata.Version,
		}
	}

	return lockfile, nil
}

func readLockFile(path, filename string) (*LockFile, bool, error) {
	loc := filepath.Join(path, filename)
	data, err := utils.AFS.ReadFile(loc)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	ld := &LockFile{
		path: path,
	}
	if err := yaml.Unmarshal(data, ld); err != nil {
		return nil, false, err
	}

	return ld, true, nil
}
