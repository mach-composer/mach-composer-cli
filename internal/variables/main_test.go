package variables

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/labd/mach-composer/internal/utils"
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

	err = vars.Load(context.Background(), "variables.yaml")
	assert.NoError(t, err)

	expected := map[string]Value{
		"foo.bar.secrets.foo": {
			val:        "encrypted",
			fileSource: &FileSource{Filename: "variables.yaml"}},
	}
	assert.EqualValues(t, expected, vars.vars)
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
