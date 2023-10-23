package lockfile

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestWriteLock(t *testing.T) {
	fileName := path.Join("./testdata/", FileName)

	err := WriteLock(&LockFile{
		path:          "./testdata",
		TerraformHash: "tf-hash",
		ConfigHash:    "cfg-hash",
		Version:       "unknown",
	})
	assert.NoError(t, err)
	assert.FileExists(t, fileName)

	b, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	s := string(b)
	assert.Equal(t, "tf-hash: tf-hash\ncfg-hash: cfg-hash\nversion: unknown\n", s)

	t.Cleanup(func() {
		err := os.Remove(fileName)
		if err != nil {
			panic(err)
		}
	})

}
