package config

import (
	"context"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func TestNewVariablesFromFile(t *testing.T) {
	utils.FS = afero.NewMemMapFs()
	utils.AFS = &afero.Afero{Fs: utils.FS}

	content := utils.TrimIndent(`
		foo:
		  bar:
			secrets:
			  foo: encrypted
	`)

	err := utils.AFS.WriteFile("variables.yaml", []byte(content), 0644)
	require.NoError(t, err)

	vars := NewVariables()
	assert.NoError(t, err)

	err = vars.Load(context.Background(), "variables.yaml", ".")
	assert.NoError(t, err)

	expected := map[string]Value{
		"foo.bar.secrets.foo": {
			val:        "encrypted",
			fileSource: &FileSource{Filename: "variables.yaml"}},
	}
	assert.EqualValues(t, expected, vars.vars)
}

func TestEncryptedVariables(t *testing.T) {
	content, err := os.ReadFile("testdata/secrets.enc.yaml")
	require.NoError(t, err)

	utils.FS = afero.NewMemMapFs()
	utils.AFS = &afero.Afero{Fs: utils.FS}

	vars := NewVariables()

	err = utils.AFS.WriteFile("testdata/secrets.enc.yaml", content, 0600)
	require.NoError(t, err)

	err = vars.Load(context.Background(), "testdata/secrets.enc.yaml", ".")
	require.NoError(t, err)

	expected := map[string]Value{
		"secrets.my-service.username": {
			val: "ENC[AES256_GCM,data:OUOm677N57JDXuEfIrk1Fhew,iv:AGMwhoqB0KwNMiDhFBZmYaIW4hoDw+75Y36+MRPaTx4=,tag:8fX4amlPMqu0kZ8uLTa6Kw==,type:str]",
			fileSource: &FileSource{
				Filename:  "testdata/secrets.enc.yaml",
				Encrypted: true,
			}},
		"secrets.my-service.password": {
			val: "ENC[AES256_GCM,data:8koAST5MJlIfao1GM4G1KTcj,iv:2XA2AqcFguEwtHTygq1KpoefkTZ2rUvlLblSjh7ZO5Y=,tag:69O8A7UIqG7QU9zxQ+0whw==,type:str]",
			fileSource: &FileSource{
				Filename:  "testdata/secrets.enc.yaml",
				Encrypted: true,
			}},
	}
	assert.EqualValues(t, expected, vars.vars)

	val, err := vars.getValue("my-site", "var.secrets.my-service.username")
	require.NoError(t, err)
	assert.Equal(t, `${data.sops_external.variables.data["secrets.my-service.username"]}`, val)

	val, err = vars.getValue("my-site", "var.secrets.my-service.password")
	require.NoError(t, err)
	assert.Equal(t, `${data.sops_external.variables.data["secrets.my-service.password"]}`, val)

	hasEncrypted := vars.HasEncrypted("my-site")
	assert.True(t, hasEncrypted)

	fs := vars.GetEncryptedSources("my-site")
	assert.Len(t, fs, 1)
}

func TestEnvVar(t *testing.T) {
	vars := NewVariables()
	t.Setenv("MY_ENV", "hello world")
	val, err := vars.getValue("my-site", "env.MY_ENV")
	require.NoError(t, err)
	assert.Equal(t, "hello world", val)
}

func TestSerializeNestedVariables(t *testing.T) {
	input := map[string]any{
		"foo": "bar",
		"level-1": map[string]any{
			"int":    10,
			"string": "my-string",
			"level-2": map[string]any{
				"int":    20,
				"string": "my-nestedstring",
			},
		},
	}
	expected := map[string]Value{
		"foo":                    {val: "bar"},
		"level-1.int":            {val: "10"},
		"level-1.string":         {val: "my-string"},
		"level-1.level-2.int":    {val: "20"},
		"level-1.level-2.string": {val: "my-nestedstring"},
	}
	result := map[string]Value{}
	serializeNestedVariables(input, result, "")
	assert.Equal(t, expected, result)
}

func TestVariablesResolve(t *testing.T) {
	data := []byte(utils.TrimIndent(`
        sites:
        - identifier: my-site
          endpoints:
            main: api.my-site.nl
            internal:
              url: internal-api.my-site.nl
              throttling_burst_limit: 5000
              throttling_rate_limit: 10000
          aws:
            account_id: 123456789
            region: eu-central-1
          components:
          - name: your-component
            variables:
              FOO_VAR: my-value
			  BAR_VAR: ${var.my-foo}
			  MULTIPLE_VARS: ${var.foo.bar} ${var.bar.foo}
            secrets:
              MY_SECRET: secretvalue
	`))
	node := yaml.Node{}
	err := yaml.Unmarshal(data, &node)
	require.NoError(t, err)

	vars := NewVariables()
	vars.vars["my-foo"] = Value{val: "my-very-special-foo"}
	vars.vars["foo.bar"] = Value{val: "my-other-bar"}
	vars.vars["bar.foo"] = Value{val: "my--bar"}

	err = vars.InterpolateNode(&node)
	require.NoError(t, err)
}
