package updater

import (
	"testing"

	"github.com/labd/mach-composer-go/config"
	"github.com/labd/mach-composer-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestMachConfigUpdater(t *testing.T) {

	data := []byte(utils.TrimIndent(`
        ---
        components:
        - name: your-component
          source: "git::https://github.com/<username>/<your-component>.git//terraform"
          version: 0.1.0
          endpoints:
            internal: internal
          integrations:
            - aws
            - commercetools
        -
          source: "git::https://github.com/<username>/<your-component>.git//terraform"
          version: "version"
          name: "foobar"
	`))

	updates := &UpdateSet{
		filename: "foo.yml",
		components: []ComponentUpdate{
			{
				component: &config.Component{
					Name:    "your-component",
					Version: "0.1.0",
				},
				version: "new-version",
			},
			{
				component: &config.Component{
					Name:    "foobar",
					Version: "version",
				},
				version: "quoted",
			},
			{
				component: &config.Component{
					Name:    "errror",
					Version: "not-in-file",
				},
				version: "quoted",
			},
		},
	}

	output := MachConfigUpdater(data, updates)
	expected := []byte(utils.TrimIndent(`
        ---
        components:
        - name: your-component
          source: "git::https://github.com/<username>/<your-component>.git//terraform"
          version: "new-version"
          endpoints:
            internal: internal
          integrations:
            - aws
            - commercetools
        -
          source: "git::https://github.com/<username>/<your-component>.git//terraform"
          version: "quoted"
		  name: "foobar"
	`))

	assert.EqualValues(t, string(expected), string(output))

}
