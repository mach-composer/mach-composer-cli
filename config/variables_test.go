package config

import (
	"testing"

	"github.com/labd/mach-composer-go/utils"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestVariablesLoad(t *testing.T) {

	utils.FS = afero.NewMemMapFs()
	utils.AFS = &afero.Afero{Fs: utils.FS}

	content := utils.TrimIndent(`
		foo:
		  bar:
			secrets:
			  foo: encrypted
		sops:
		  kms: []
		  gcp_kms: []
		  azure_kv:
			- vault_url: https://gdsharedwekvsops.vault.azure.net
			  name: sops-key
			  version: aa5a053752be4cee83ee0d57200fb1a8
			  created_at: "2021-06-03T21:15:12Z"
			  enc: secret
		  hc_vault: []
		  age: []
		  lastmodified: "2022-03-21T10:44:39Z"
		  mac: mymac
		  pgp: []
		  unencrypted_suffix: _unencrypted
		  version: 3.7.1
	`)

	utils.AFS.WriteFile("variables.yaml", []byte(content), 0644)

	vars, err := loadVariables("variables.yaml")
	assert.NoError(t, err)

	expected := map[string]string{
		"foo.bar.secrets.foo": "encrypted",
	}
	assert.EqualValues(t, expected, vars.vars)
}
