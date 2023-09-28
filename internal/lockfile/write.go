package lockfile

import (
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"gopkg.in/yaml.v3"
	"path/filepath"
)

func WriteLock(lf *LockFile) error {
	data, err := yaml.Marshal(lf)
	if err != nil {
		return err
	}

	filename := filepath.Join(lf.path, FileName)
	if err := utils.AFS.WriteFile(filename, data, 0600); err != nil {
		return err
	}

	return nil
}
