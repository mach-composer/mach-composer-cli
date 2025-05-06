package state

import (
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

// HttpState HTTP backend configuration.
type HttpState struct {
	Address                string `mapstructure:"address"`
	UpdateMethod           string `mapstructure:"update_method" default:"POST"`
	LockAddress            string `mapstructure:"lock_address"`
	LockMethod             string `mapstructure:"lock_method" default:"LOCK"`
	UnlockAddress          string `mapstructure:"unlock_address"`
	UnlockMethod           string `mapstructure:"unlock_method" default:"UNLOCK"`
	Username               string `mapstructure:"username"`
	Password               string `mapstructure:"password"`
	SkipCertVerification   bool   `mapstructure:"skip_cert_verification" default:"false"`
	RetryMax               int    `mapstructure:"retry_max" default:"2"`
	RetryWaitMin           int    `mapstructure:"retry_wait_min" default:"1"`
	RetryWaitMax           int    `mapstructure:"retry_wait_max" default:"30"`
	ClientCertificatePEM   string `mapstructure:"client_certificate_pem"`
	ClientPrivateKeyPEM    string `mapstructure:"client_private_key_pem"`
	ClientCACertificatePEM string `mapstructure:"client_ca_certificate_pem"`
}

type HttpRenderer struct {
	BaseRenderer
	state *HttpState
}

func (hr *HttpRenderer) Backend() (string, error) {
	templateContext := struct {
		State      *HttpState
		Identifier string
	}{
		State:      hr.state,
		Identifier: hr.identifier,
	}

	tpl := `
	backend "http" {
		address               = "{{ .State.Address }}"
		update_method         = "{{ .State.UpdateMethod }}"
		{{- if .State.LockAddress }}
		lock_address          = "{{ .State.LockAddress }}"
		{{- end }}
		{{- if .State.LockMethod }}
		lock_method           = "{{ .State.LockMethod }}"
		{{- end }}
		{{- if .State.UnlockAddress }}
		unlock_address        = "{{ .State.UnlockAddress }}"
		{{- end }}
		{{- if .State.UnlockMethod }}
		unlock_method         = "{{ .State.UnlockMethod }}"
		{{- end }}
		{{- if .State.Username }}
		username              = "{{ .State.Username }}"
		{{- end }}
		{{- if .State.Password }}
		password              = "{{ .State.Password }}"
		{{- end }}
		skip_cert_verification = {{ .State.SkipCertVerification }}
		{{- if gt .State.RetryMax 0 }}
		retry_max             = {{ .State.RetryMax }}
		{{- end }}
		{{- if gt .State.RetryWaitMin 0 }}
		retry_wait_min        = {{ .State.RetryWaitMin }}
		{{- end }}
		{{- if gt .State.RetryWaitMax 0 }}
		retry_wait_max        = {{ .State.RetryWaitMax }}
		{{- end }}
		{{- if .State.ClientCertificatePEM }}
		client_certificate_pem = <<EOT
{{ .State.ClientCertificatePEM }}
EOT
		{{- end }}
		{{- if .State.ClientPrivateKeyPEM }}
		client_private_key_pem = <<EOT
{{ .State.ClientPrivateKeyPEM }}
EOT
		{{- end }}
		{{- if .State.ClientCACertificatePEM }}
		client_ca_certificate_pem = <<EOT
{{ .State.ClientCACertificatePEM }}
EOT
		{{- end }}
	}
	`
	return utils.RenderGoTemplate(tpl, templateContext)
}

func (hr *HttpRenderer) RemoteState() (string, error) {
	templateContext := struct {
		State      *HttpState
		Identifier string
		Key        string
	}{
		State:      hr.state,
		Identifier: hr.identifier,
		Key:        hr.stateKey,
	}

	tpl := `
	data "terraform_remote_state" "{{ .Key }}" {
		backend = "http"

		config = {
			address               = "{{ .State.Address }}"
			update_method         = "{{ .State.UpdateMethod }}"
			{{- if .State.Username }}
			username              = "{{ .State.Username }}"
			{{- end }}
			{{- if .State.Password }}
			password              = "{{ .State.Password }}"
			{{- end }}
			skip_cert_verification = {{ .State.SkipCertVerification }}
			{{- if gt .State.RetryMax 0 }}
			retry_max             = {{ .State.RetryMax }}
			{{- end }}
			{{- if gt .State.RetryWaitMin 0 }}
			retry_wait_min        = {{ .State.RetryWaitMin }}
			{{- end }}
			{{- if gt .State.RetryWaitMax 0 }}
			retry_wait_max        = {{ .State.RetryWaitMax }}
			{{- end }}
		}
	}
	`
	return utils.RenderGoTemplate(tpl, templateContext)
}
