package lockfile

import (
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadLockFileNotExist(t *testing.T) {
	_, exists, err := readLockFile("./testdata", "i-dont-exist.txt")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestReadLockFileErr(t *testing.T) {
	_, exists, err := readLockFile(".", "")
	assert.Error(t, err)
	assert.False(t, exists)
}

func TestReadLockFileInvalidContent(t *testing.T) {
	_, exists, err := readLockFile("./testdata", "faulty.yaml")
	assert.Error(t, err)
	assert.False(t, exists)
}

func TestReadLockFileOk(t *testing.T) {
	lf, exists, err := readLockFile("./testdata", "ok.yaml")
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, "de17e51f77a42e160b9e1fb41e0514cf17a32be58f42466f0bc5ff6e10eda4a1", lf.ConfigHash)
	assert.Equal(t, "51fbd31ffe2b9a53497d2300101af80d724b03f3b270315f34e8efa71abc4ebd", lf.TerraformHash)
	assert.Equal(t, "unknown", lf.Version)
}

func TestGetNewLock(t *testing.T) {
	cfg := &config.MachConfig{
		ConfigHash: "hash",
	}
	lf, err := GetLock(cfg, "./testdata")
	assert.NoError(t, err)
	assert.True(t, lf.isNew)
	assert.Equal(t, "hash", lf.ConfigHash)
	assert.Equal(t, "unknown", lf.Version)
}
