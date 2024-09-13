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

func (s *GenerateTestSuite) TestGenerateAWSMultisite() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/generate/aws-multisite")

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

func (s *GenerateTestSuite) TestGenerateAWSDeploymentTypeSite() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/generate/aws-deployment-type-site")

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

func (s *GenerateTestSuite) TestGenerateAWSDeploymentTypeMixed() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/generate/aws-deployment-type-mixed")

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

func (s *GenerateTestSuite) TestGenerateAWSBasic() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/generate/aws-basic")

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

func (s *GenerateTestSuite) TestGenerateGCPBasic() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/generate/gcp-basic")

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

func (s *GenerateTestSuite) TestGenerateAzureBasic() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/generate/azure-basic")

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

func (s *GenerateTestSuite) TestGenerateLocalBasic() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/generate/local-basic")

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

func (s *GenerateTestSuite) TestGenerateWithAlias() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/generate/with-alias")

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

func (s *GenerateTestSuite) TestGenerateWithVariables() {
	pwd, _ := os.Getwd()
	workdir := path.Join(pwd, "testdata/cases/generate/with-variables")

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
