package cmd

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"os/exec"
	"path"
	"testing"
)

type ValidateTestSuite struct {
	suite.Suite
	tempDir string
}

func TestValidateTestSuite(t *testing.T) {
	suite.Run(t, new(ValidateTestSuite))
}

func (s *ValidateTestSuite) SetupSuite() {
	_, err := exec.LookPath("terraform")
	if err != nil {
		s.T().Fatal("terraform command not found")
	}

	tmpDir, _ := os.MkdirTemp("mach-composer", "test")
	_ = os.Setenv("TF_PLUGIN_CACHE_DIR", tmpDir)
	_ = os.Setenv("TF_PLUGIN_CACHE_MAY_BREAK_DEPENDENCY_LOCK_FILE", "1")

	s.tempDir = tmpDir
}

func (s *ValidateTestSuite) TearDownSuite() {
	_ = os.RemoveAll(s.tempDir)
}

func (s *ValidateTestSuite) TestValidateInvalid() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/validate/invalid")

	cmd := RootCmd
	cmd.SetArgs([]string{
		"validate",
		"--file", path.Join(workdir, "main.yaml"),
		"--validation-path", path.Join(workdir, "validations"),
	})
	err := cmd.Execute()
	assert.Error(s.T(), err)
}

func (s *ValidateTestSuite) TestValidateValid() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/validate/valid")

	cmd := RootCmd
	cmd.SetArgs([]string{
		"validate",
		"--file", path.Join(workdir, "main.yaml"),
		"--validation-path", path.Join(workdir, "validations"),
	})
	err := cmd.Execute()
	assert.NoError(s.T(), err)
}
