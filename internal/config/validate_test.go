package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/mach-composer/mach-composer-cli/internal/plugins"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func TestCreateSchema(t *testing.T) {
	data := []byte(utils.TrimIndent(`
        ---
        mach_composer:
          version: 1.0.0
		  plugins: {}
        global:
          environment: test
		  terraform_config: {}
		  cloud: aws
        sites:
        - identifier: my-site
		  aws:
		    stringValue: 027113887083
			boolValue: -1
			intValue: no-int
			mapOfInts:
		    	1324: 027113887083
          components:
          - name: your-component
            variables:
              FOO_VAR: my-value
			  BAR_VAR: ${var.foo}
			  MULTIPLE_VARS: ${var.foo.bar} ${var.bar.foo}
			  NESTED:
			  	VARS: OK
            secrets:
              MY_SECRET: secretvalue
        components:
        - name: your-component
          source: "git::https://github.com/<username>/<your-component>.git//terraform"
          version: 0.1.0
    `))

	document := &yaml.Node{}
	err := yaml.Unmarshal(data, document)
	require.NoError(t, err)

	pr := plugins.NewPluginRepository()
	err = pr.Add("aws", plugins.NewPluginV1Adapter(plugins.NewMockPluginV1()))
	require.NoError(t, err)

	raw, err := newRawConfig("main.yml", document)
	raw.plugins = pr
	require.NoError(t, err)

	isValid, err := validateCompleteConfig(raw)
	require.Error(t, err)
	assert.False(t, isValid)

	assert.ErrorContains(t, err, "sites.0.aws: requiredValue is required")
	assert.ErrorContains(t, err, "sites.0.aws.boolValue: Invalid type. Expected: boolean, given: integer")
	assert.ErrorContains(t, err, "sites.0.aws.stringValue: Invalid type. Expected: string, given: integer")
	assert.ErrorContains(t, err, "sites.0.aws.intValue: Invalid type. Expected: number, given: string")
}
