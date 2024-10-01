package gitutils

import (
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPathFilter(t *testing.T) {
	assert.True(t, pathFilter([]string{})("test/test"))
	assert.True(t, pathFilter([]string{"test"})("test/test"))
	assert.False(t, pathFilter([]string{"foo"})("test/test"))
}

func TestCommitsBetweenResolveFirstError(t *testing.T) {
	tr := NewTestRepository("")

	// First commit
	err := tr.addTextFile("test-1.txt", "test")
	require.NoError(t, err)

	// Second commit
	err = tr.addTextFile("test-2.txt", "test")
	require.NoError(t, err)

	secondHash, err := tr.commit("Second commit")
	require.NoError(t, err)

	targetRev := plumbing.Revision(secondHash.String())

	badRev := plumbing.Revision("test")

	// Check results
	_, err = commitsBetween(tr.repository(), &badRev, &targetRev, []string{})
	require.Error(t, err)
}

func TestCommitsBetween(t *testing.T) {
	tr := NewTestRepository("")

	// First commit
	err := tr.addTextFile("test-1.txt", "test")
	require.NoError(t, err)

	firstHash, err := tr.commit("Initial commit")
	require.NoError(t, err)

	// Second commit
	err = tr.addTextFile("test-2.txt", "test")
	require.NoError(t, err)

	secondHash, err := tr.commit("Second commit")
	require.NoError(t, err)

	targetRev := plumbing.Revision(secondHash.String())
	baseRev := plumbing.Revision(firstHash.String())

	commits, err := commitsBetween(tr.repository(), nil, &targetRev, []string{})
	require.NoError(t, err)
	assert.Equal(t, 2, len(commits))

	commits, err = commitsBetween(tr.repository(), &baseRev, &targetRev, []string{})
	require.NoError(t, err)
	assert.Equal(t, 1, len(commits))
}

func TestCommitsBetweenFilterPath(t *testing.T) {
	tr := NewTestRepository("")

	// First commit
	err := tr.addTextFile("test-1.txt", "test")
	require.NoError(t, err)

	firstHash, err := tr.commit("Initial commit")
	require.NoError(t, err)
	require.NotNil(t, firstHash)

	// Second commit
	err = tr.addTextFile("test-2.txt", "test")
	require.NoError(t, err)

	err = tr.addTextFile("wanted/test-2.txt", "test")
	require.NoError(t, err)

	_, err = tr.commit("Second commit")
	require.NoError(t, err)

	// Third commit
	err = tr.addTextFile("test-3.txt", "test")
	require.NoError(t, err)

	err = tr.addTextFile("wantedotherdir/test-3.txt", "test")
	require.NoError(t, err)

	thirdHash, err := tr.commit("third commit")
	require.NoError(t, err)

	targetRev := plumbing.Revision(thirdHash.String())

	// Check results
	commits, err := commitsBetween(tr.repository(), nil, &targetRev, []string{})
	require.NoError(t, err)
	assert.Equal(t, 3, len(commits))

	commits, err = commitsBetween(tr.repository(), nil, &targetRev, []string{"wanted/"})
	require.NoError(t, err)
	assert.Equal(t, 1, len(commits))
}
