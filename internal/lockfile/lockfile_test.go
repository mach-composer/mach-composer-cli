package lockfile

import (
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasChangesNew(t *testing.T) {
	lf := &LockFile{
		isNew: true,
	}

	cfg := &config.MachConfig{}

	assert.True(t, lf.HasChanges(cfg))
}

func TestHasChangesDifferentConfigHash(t *testing.T) {
	lf := &LockFile{
		ConfigHash: "hash",
	}

	cfg := &config.MachConfig{
		ConfigHash: "different-hash",
	}

	assert.True(t, lf.HasChanges(cfg))
}

func TestHasChangesDifferentVersion(t *testing.T) {
	lf := &LockFile{
		ConfigHash: "hash",
		Version:    "some-version",
	}

	cfg := &config.MachConfig{
		ConfigHash: "hash",
	}

	assert.True(t, lf.HasChanges(cfg))
}

func TestHasChangesNoChanges(t *testing.T) {
	lf := &LockFile{
		ConfigHash: "hash",
		Version:    "unknown",
	}

	cfg := &config.MachConfig{
		ConfigHash: "hash",
	}

	assert.False(t, lf.HasChanges(cfg))
}

func TestUpdate(t *testing.T) {
	lf := &LockFile{
		path:          "./testdata",
		TerraformHash: "tf-hash",
		ConfigHash:    "cfg-hash",
		Version:       "v1",
	}

	err := lf.Update(&config.MachConfig{
		ConfigHash: "new-cfg-hash",
	})
	assert.NoError(t, err)

	assert.Equal(t, "new-cfg-hash", lf.ConfigHash)
	assert.NotEqual(t, "tf-hash", lf.TerraformHash)
	assert.Equal(t, "unknown", lf.Version)
}
