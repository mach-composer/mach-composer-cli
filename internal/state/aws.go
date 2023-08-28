package state

import (
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

// AwsState AWS S3 bucket state backend configuration.
type AwsState struct {
	Bucket    string `mapstructure:"bucket"`
	KeyPrefix string `mapstructure:"key_prefix"`
	Region    string `mapstructure:"region"`
	RoleARN   string `mapstructure:"role_arn"`
	LockTable string `mapstructure:"lock_table"`
	Encrypt   bool   `mapstructure:"encrypt" default:"true"`
}

type AwsRenderer struct {
	key   string
	state *AwsState
}

func (ar *AwsRenderer) Key() string {
	return ar.key
}

func (ar *AwsRenderer) Backend() (string, error) {
	templateContext := struct {
		State *AwsState
		Key   string
	}{
		State: ar.state,
		Key:   ar.key,
	}

	tpl := `
	backend "s3" {
	  bucket         = "{{ .State.Bucket }}"
	  key            = "{{ .State.KeyPrefix}}/{{ .Key }}"
	  region         = "{{ .State.Region }}"
	  {{ if .State.RoleARN }}
	  role_arn       = "{{ .State.RoleARN }}"
	  {{ end }}
	  {{ if .State.LockTable }}
	  dynamodb_table = "{{ .State.LockTable }}"
	  {{ end }}
	  encrypt        = {{ .State.Encrypt }}
	}
	`
	return utils.RenderGoTemplate(tpl, templateContext)
}

func (ar *AwsRenderer) RemoteState() (string, error) {
	templateContext := struct {
		State *AwsState
		Key   string
	}{
		State: ar.state,
		Key:   ar.key,
	}

	template := `
	data "terraform_remote_state" "{{ .Key }}" {
	  backend = "aws"
	
	  config = {
		  bucket         = "{{ .State.Bucket }}"
		  key            = "{{ .State.KeyPrefix}}/{{ .Key }}"
		  region         = "{{ .State.Region }}"
		  {{ if .State.RoleARN }}
		  role_arn       = "{{ .State.RoleARN }}"
		  {{ end }}
		  {{ if .State.LockTable }}
		  dynamodb_table = "{{ .State.LockTable }}"
		  {{ end }}
		  encrypt        = {{ .State.Encrypt }}
	  }
	}
	`
	return utils.RenderGoTemplate(template, templateContext)
}
