package lockfile

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasChangesNew(t *testing.T) {
	lf := &LockFile{
		isNew: true,
	}

	assert.True(t, lf.HasChanges(""))
}

func TestHasChangesDifferentConfigHash(t *testing.T) {
	lf := &LockFile{
		ConfigHash: "hash",
	}
	assert.True(t, lf.HasChanges("different-hash"))
}

func TestHasChangesDifferentVersion(t *testing.T) {
	lf := &LockFile{
		ConfigHash: "hash",
		Version:    "some-version",
	}

	assert.True(t, lf.HasChanges("hash"))
}

func TestHasChangesNoChanges(t *testing.T) {
	lf := &LockFile{
		ConfigHash: "hash",
		Version:    "unknown",
	}
	assert.False(t, lf.HasChanges("hash"))
}

func TestUpdate(t *testing.T) {
	lf := &LockFile{
		path:          "./testdata",
		TerraformHash: "tf-hash",
		ConfigHash:    "cfg-hash",
		Version:       "v1",
	}

	err := lf.Update("new-cfg-hash")
	assert.NoError(t, err)

	assert.Equal(t, "new-cfg-hash", lf.ConfigHash)
	assert.NotEqual(t, "tf-hash", lf.TerraformHash)
	assert.Equal(t, "unknown", lf.Version)
}
