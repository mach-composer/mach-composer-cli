package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ComputeDirHash_Same(t *testing.T) {
	h1, err := ComputeDirHash("testdata/dirhash")
	assert.NoError(t, err)

	h2, err := ComputeDirHash("testdata/dirhash")
	assert.NoError(t, err)

	assert.Equal(t, h1, h2)
}

func Test_ComputeDirHash_Different(t *testing.T) {
	h1, err := ComputeDirHash("testdata/dirhash")
	assert.NoError(t, err)

	h2, err := ComputeDirHash("testdata/dirhash2")
	assert.NoError(t, err)

	assert.NotEqual(t, h1, h2)
}
