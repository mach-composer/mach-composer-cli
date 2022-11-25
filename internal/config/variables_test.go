package config

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/labd/mach-composer/internal/utils"
)

func TestVariablesLoad(t *testing.T) {
	utils.FS = afero.NewMemMapFs()
	utils.AFS = &afero.Afero{Fs: utils.FS}

	content := utils.TrimIndent(`
		foo:
		  bar:
			secrets:
			  foo: encrypted
	`)

	utils.AFS.WriteFile("variables.yaml", []byte(content), 0644)

	vars, err := loadVariables(context.Background(), "variables.yaml")
	assert.NoError(t, err)

	expected := map[string]string{
		"foo.bar.secrets.foo": "encrypted",
	}
	assert.EqualValues(t, expected, vars.vars)
}
