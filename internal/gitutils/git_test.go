package gitutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGitSource(t *testing.T) {
	params := []struct {
		source   string
		expected *gitSource
	}{
		{
			source: "git::https://github.com/labd/mach-composer",
			expected: &gitSource{
				Repository: "https://github.com/labd/mach-composer",
				Path:       "",
				Name:       "mach-composer",
			},
		},
		{
			source: "git::https://github.com/labd/mach-composer//terraform",
			expected: &gitSource{
				Repository: "https://github.com/labd/mach-composer",
				Path:       "terraform",
				Name:       "mach-composer",
			},
		},
	}
	for _, p := range params {
		p.expected.URL = p.source

		res, err := parseGitSource(p.source)
		assert.NoError(t, err)
		assert.EqualValues(t, p.expected, res)
	}
}
