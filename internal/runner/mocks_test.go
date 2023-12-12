package runner

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type TerraformRunnerMock struct {
	mock.Mock
}

func (m *TerraformRunnerMock) RunTerraform(ctx context.Context, cwd string, args ...string) error {
	argv := m.Called(ctx, cwd, args)
	return argv.Error(0)
}
