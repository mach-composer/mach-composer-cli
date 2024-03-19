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
		{source: Source("http://example.com/vpc.zip"), expected: SourceTypeHttp},
		{source: Source("https://example.com/vpc.zip"), expected: SourceTypeHttp},
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
		{source: Source("git::https://example.com/vpc.git"), expected: "git::https://example.com/vpc.git?ref=v1.0.0"},
		{source: Source("s3::https://github.com/my/project"), expected: "s3::https://github.com/my/project/v1.0.0.zip"},
		{source: Source("gcs::https://github.com/my/project"), expected: "gcs::https://github.com/my/project/v1.0.0.zip"},

		//Unsupported sources. We throw an error instead, as we don't know how to handle them yet
		{source: Source("http://example.com/vpc.zip"), error: true},
		{source: Source("https://example.com/vpc.zip"), error: true},
		{source: Source("git@github.com:hashicorp/example.git"), error: true},
		{source: Source("github.com/hashicorp/example"), error: true},
		{source: Source("bitbucket.org/hashicorp/terraform-consul-aws"), error: true},
		{source: Source("hg::http://example.com/vpc.hg"), error: true},
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
