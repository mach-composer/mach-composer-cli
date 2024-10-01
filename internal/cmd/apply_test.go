//go:build integration

package cmd

import (
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"path"
	"testing"
)

import (
	"github.com/stretchr/testify/suite"
)

type ApplyTestSuite struct {
	suite.Suite
	tempDir string
}

func TestApplyTestSuite(t *testing.T) {
	t.Skip("skipping test")
	suite.Run(t, new(ApplyTestSuite))
}

func (s *ApplyTestSuite) SetupSuite() {
	_, err := exec.LookPath("terraform")
	if err != nil {
		s.T().Fatal("terraform command not found")
	}

	tmpDir, _ := os.MkdirTemp("mach-composer", "test")
	_ = os.Setenv("TF_PLUGIN_CACHE_DIR", tmpDir)
	_ = os.Setenv("TF_PLUGIN_CACHE_MAY_BREAK_DEPENDENCY_LOCK_FILE", "1")

	s.tempDir = tmpDir
}

func (s *ApplyTestSuite) TearDownSuite() {
	_ = os.RemoveAll(s.tempDir)
}

func (s *ApplyTestSuite) TestApplySimple() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/apply/simple")
	defer cleanWorkingDir(workdir)

	cmd := RootCmd
	_ = os.Setenv("MC_HASH_FILE", path.Join(workdir, "hashes.json"))
	cmd.SetArgs([]string{
		"apply",
		"--output-path", path.Join(workdir, "deployments"),
		"--file", path.Join(workdir, "main.yaml"),
		"--auto-approve",
	})
	err := cmd.Execute()
	assert.NoError(s.T(), err)

	assert.FileExists(s.T(), path.Join(workdir, "hashes.json"))
	assert.FileExists(s.T(), path.Join(workdir, "deployments/main/test-1/main.tf"))
	assert.FileExists(s.T(), path.Join(workdir, "deployments/main/test-1/states/test-1.tfstate"))
	assert.FileExists(s.T(), path.Join(workdir, "deployments/main/test-1/outputs/component-1.json"))
}

func (s *ApplyTestSuite) TestApplySplitState() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/apply/split-state")
	//defer cleanWorkingDir(workdir)

	cmd := RootCmd
	_ = os.Setenv("MC_HASH_FILE", path.Join(workdir, "deployments/hashes.json"))
	_ = os.Setenv("STATES_PATH", path.Join(workdir, "states"))
	cmd.SetArgs([]string{
		"apply",
		"--output-path", path.Join(workdir, "deployments"),
		"--file", path.Join(workdir, "main.yaml"),
		"--auto-approve",
	})
	err := cmd.Execute()
	assert.NoError(s.T(), err)

	assert.FileExists(s.T(), path.Join(workdir, "deployments/hashes.json"))
	assert.FileExists(s.T(), path.Join(workdir, "deployments/main/test-1/main.tf"))
	assert.FileExists(s.T(), path.Join(workdir, "deployments/main/test-1/component-2/main.tf"))
	assert.FileExists(s.T(), path.Join(workdir, "states/test-1.tfstate"))
	assert.FileExists(s.T(), path.Join(workdir, "states/test-1/component-2.tfstate"))
}

func (s *ApplyTestSuite) TestApplyNoHashesFile() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/apply/simple")
	defer cleanWorkingDir(workdir)

	cmd := RootCmd
	_ = os.Setenv("MC_HASH_FILE", path.Join(workdir, "hashes.json"))
	cmd.SetArgs([]string{
		"apply",
		"--output-path", path.Join(workdir, "deployments"),
		"--file", path.Join(workdir, "main.yaml"),
		"--auto-approve",
	})
	err := cmd.Execute()
	assert.NoError(s.T(), err)

	assert.FileExists(s.T(), path.Join(workdir, "hashes.json"))

	err = os.RemoveAll(path.Join(workdir, "hashes.json"))
	if err != nil {
		s.T().Fatal(err)
	}

	err = cmd.Execute()
	assert.NoError(s.T(), err)
	assert.FileExists(s.T(), path.Join(workdir, "hashes.json"))
}
