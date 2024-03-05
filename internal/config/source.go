package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

type SourceType string

const (
	SourceTypeLocal     SourceType = "local"
	SourceTypeGit       SourceType = "git"
	SourceTypeGithub    SourceType = "github"
	SourceTypeBitbucket SourceType = "bitbucket"
	SourceTypeMercurial SourceType = "mercurial"
	SourceTypeHttp      SourceType = "http"
	SourceTypeS3        SourceType = "s3"
	SourceTypeGCS       SourceType = "gcs"
)

type Source string

func (s *Source) String() string {
	return string(*s)
}

func (s *Source) IsType(t SourceType) bool {
	st, err := s.Type()
	if err != nil {
		return false
	}

	return st == t
}

func (s *Source) Type() (SourceType, error) {

	if strings.HasPrefix(string(*s), ".") {
		return SourceTypeLocal, nil
	}

	if strings.HasPrefix(string(*s), "github.com") ||
		strings.HasPrefix(string(*s), "git@") {
		return SourceTypeGithub, nil
	}

	if strings.HasPrefix(string(*s), "bitbucket.org") {
		return SourceTypeBitbucket, nil
	}

	if strings.HasPrefix(string(*s), "http") {
		return SourceTypeHttp, nil
	}

	if strings.HasPrefix(string(*s), "git::") {
		return SourceTypeGit, nil
	}

	if strings.HasPrefix(string(*s), "hg::") {
		return SourceTypeMercurial, nil
	}

	if strings.HasPrefix(string(*s), "s3::") {
		return SourceTypeS3, nil
	}

	if strings.HasPrefix(string(*s), "gcs::") {
		return SourceTypeGCS, nil
	}

	return "", fmt.Errorf("unknown source type: %s", string(*s))
}

func (s *Source) GetVersionSource(version string) (string, error) {
	t, err := s.Type()
	if err != nil {
		return "", err
	}

	switch t {
	case SourceTypeLocal:
		// If the source is a relative locale path then transform it to an
		// absolute path (required for Terraform)
		return filepath.Abs(string(*s))
	case SourceTypeGit:
		// When using Git, we will automatically add a reference to the string
		// so that the given version is used when fetching the module itself
		// from Git as well
		return fmt.Sprintf("%s?ref=%s", string(*s), version), nil
	case SourceTypeS3:
		fallthrough
	case SourceTypeGCS:
		// For GCS and AWS we assume that the version is the name of the zip file
		return fmt.Sprintf("%s/%s.zip", string(*s), version), nil
	default:
		// For all other sources we will just return the source as is
		return string(*s), nil
	}
}
