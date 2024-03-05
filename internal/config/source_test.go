package config

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestType(t *testing.T) {
	tests := []struct {
		source   Source
		expected SourceType
		error    bool
	}{
		{source: Source("./local/file"), expected: SourceTypeLocal},
		{source: Source("../local/file"), expected: SourceTypeLocal},
		{source: Source("git@github.com:hashicorp/example.git"), expected: SourceTypeGithub},
		{source: Source("github.com/hashicorp/example"), expected: SourceTypeGithub},
		{source: Source("bitbucket.org/hashicorp/terraform-consul-aws"), expected: SourceTypeBitbucket},
		{source: Source("git::https://example.com/vpc.git"), expected: SourceTypeGit},
		{source: Source("hg::http://example.com/vpc.hg"), expected: SourceTypeMercurial},
		{source: Source("s3::https://github.com/my/project"), expected: SourceTypeS3},
		{source: Source("gcs::https://github.com/my/project"), expected: SourceTypeGCS},
		{source: Source("bla::https://github.com/my/project"), error: true},
	}
	for _, tc := range tests {
		t.Run(string(tc.source), func(t *testing.T) {
			result, err := tc.source.Type()
			if tc.error {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func mustAbs(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		panic(err)
	}
	return abs
}

func TestGetVersionSource(t *testing.T) {
	tests := []struct {
		source   Source
		expected string
		error    bool
	}{
		{source: Source("./local/file"), expected: mustAbs("./local/file")},
		{source: Source("../local/file"), expected: mustAbs("../local/file")},
		{source: Source("git@github.com:hashicorp/example.git"), expected: "git@github.com:hashicorp/example.git"},
		{source: Source("github.com/hashicorp/example"), expected: "github.com/hashicorp/example"},
		{source: Source("bitbucket.org/hashicorp/terraform-consul-aws"), expected: "bitbucket.org/hashicorp/terraform-consul-aws"},
		{source: Source("git::https://example.com/vpc.git"), expected: "git::https://example.com/vpc.git?ref=v1.0.0"},
		{source: Source("hg::http://example.com/vpc.hg"), expected: "hg::http://example.com/vpc.hg"},
		{source: Source("s3::https://github.com/my/project"), expected: "s3::https://github.com/my/project/v1.0.0.zip"},
		{source: Source("gcs::https://github.com/my/project"), expected: "gcs::https://github.com/my/project/v1.0.0.zip"},
		{source: Source("bla::https://github.com/my/project"), error: true},
	}
	for _, tc := range tests {
		t.Run(string(tc.source), func(t *testing.T) {
			result, err := tc.source.GetVersionSource("v1.0.0")
			if tc.error {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsTypeOK(t *testing.T) {
	s := Source("git::https://example.com/vpc.git")
	assert.True(t, s.IsType(SourceTypeGit))
}

func TestIsTypeInvalid(t *testing.T) {
	s := Source("git::https://example.com/vpc.git")
	assert.False(t, s.IsType(SourceTypeLocal))
}
