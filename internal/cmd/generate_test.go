package cmd

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"os/exec"
	"path"
	"testing"
)

type GenerateTestSuite struct {
	suite.Suite
	tempDir string
}

func TestGenerateTestSuite(t *testing.T) {
	suite.Run(t, new(GenerateTestSuite))
}

func (s *GenerateTestSuite) SetupSuite() {
	_, err := exec.LookPath("terraform")
	if err != nil {
		s.T().Fatal("terraform command not found")
	}

	tmpDir, _ := os.MkdirTemp("mach-composer", "test")
	_ = os.Setenv("TF_PLUGIN_CACHE_DIR", tmpDir)
	_ = os.Setenv("TF_PLUGIN_CACHE_MAY_BREAK_DEPENDENCY_LOCK_FILE", "1")

	s.tempDir = tmpDir
}

func (s *GenerateTestSuite) TearDownSuite() {
	_ = os.RemoveAll(s.tempDir)
}

func (s *GenerateTestSuite) TestGenerateLocalMultisite() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/generate/local-multisite")
	defer cleanWorkingDir(workdir)

	cmd := RootCmd
	cmd.SetArgs([]string{
		"generate",
		"--output-path", path.Join(workdir, "deployments"),
		"--file", path.Join(workdir, "main.yaml"),
	})
	err := cmd.Execute()
	assert.NoError(s.T(), err)
	assert.NoError(s.T(), CompareDirectories(path.Join(workdir, "deployments"), path.Join(workdir, "expected")))
}

func (s *GenerateTestSuite) TestGenerateLocalDeploymentTypeSite() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/generate/local-deployment-type-site")
	defer cleanWorkingDir(workdir)

	cmd := RootCmd
	cmd.SetArgs([]string{
		"generate",
		"--output-path", path.Join(workdir, "deployments"),
		"--file", path.Join(workdir, "main.yaml"),
	})
	err := cmd.Execute()
	assert.NoError(s.T(), err)
	assert.NoError(s.T(), CompareDirectories(path.Join(workdir, "deployments"), path.Join(workdir, "expected")))
}

func (s *GenerateTestSuite) TestGenerateLocalDeploymentTypeMixed() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/generate/local-deployment-type-mixed")
	defer cleanWorkingDir(workdir)

	cmd := RootCmd
	cmd.SetArgs([]string{
		"generate",
		"--output-path", path.Join(workdir, "deployments"),
		"--file", path.Join(workdir, "main.yaml"),
	})
	err := cmd.Execute()
	assert.NoError(s.T(), err)
	assert.NoError(s.T(), CompareDirectories(path.Join(workdir, "deployments"), path.Join(workdir, "expected")))
}

func (s *GenerateTestSuite) TestGenerateAWS() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/generate/aws-basic")
	defer cleanWorkingDir(workdir)

	cmd := RootCmd
	cmd.SetArgs([]string{
		"generate",
		"--output-path", path.Join(workdir, "deployments"),
		"--file", path.Join(workdir, "main.yaml"),
	})
	err := cmd.Execute()
	assert.NoError(s.T(), err)
	assert.NoError(s.T(), CompareDirectories(path.Join(workdir, "deployments"), path.Join(workdir, "expected")))
}

func (s *GenerateTestSuite) TestGenerateGCP() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/generate/gcp-basic")
	defer cleanWorkingDir(workdir)

	cmd := RootCmd
	cmd.SetArgs([]string{
		"generate",
		"--output-path", path.Join(workdir, "deployments"),
		"--file", path.Join(workdir, "main.yaml"),
	})
	err := cmd.Execute()
	assert.NoError(s.T(), err)
	assert.NoError(s.T(), CompareDirectories(path.Join(workdir, "deployments"), path.Join(workdir, "expected")))
}

func (s *GenerateTestSuite) TestGenerateAzure() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/generate/azure-basic")
	defer cleanWorkingDir(workdir)

	cmd := RootCmd
	cmd.SetArgs([]string{
		"generate",
		"--output-path", path.Join(workdir, "deployments"),
		"--file", path.Join(workdir, "main.yaml"),
	})
	err := cmd.Execute()
	assert.NoError(s.T(), err)
	assert.NoError(s.T(), CompareDirectories(path.Join(workdir, "deployments"), path.Join(workdir, "expected")))
}
