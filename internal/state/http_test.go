package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHttpRendererBackend(t *testing.T) {
	r := HttpRenderer{
		state: &HttpState{
			Address:                "https://example.com/state",
			UpdateMethod:           "PUT",
			LockAddress:            "https://example.com/lock",
			LockMethod:             "LOCK",
			UnlockAddress:          "https://example.com/unlock",
			UnlockMethod:           "UNLOCK",
			Username:               "user",
			Password:               "pass",
			SkipCertVerification:   true,
			RetryMax:               3,
			RetryWaitMin:           2,
			RetryWaitMax:           10,
			ClientCertificatePEM:   "client-cert",
			ClientPrivateKeyPEM:    "client-key",
			ClientCACertificatePEM: "ca-cert",
		},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
		},
	}

	b, err := r.Backend()
	assert.NoError(t, err)
	assert.Equal(t, `
	backend "http" {
		address               = "https://example.com/state"
		update_method         = "PUT"
		lock_address          = "https://example.com/lock"
		lock_method           = "LOCK"
		unlock_address        = "https://example.com/unlock"
		unlock_method         = "UNLOCK"
		username              = "user"
		password              = "pass"
		skip_cert_verification = true
		retry_max             = 3
		retry_wait_min        = 2
		retry_wait_max        = 10
		client_certificate_pem = <<EOT
client-cert
EOT
		client_private_key_pem = <<EOT
client-key
EOT
		client_ca_certificate_pem = <<EOT
ca-cert
EOT
	}
	`, b)
}

func TestHttpRendererRemoteState(t *testing.T) {
	r := HttpRenderer{
		state: &HttpState{
			Address:              "https://example.com/state",
			UpdateMethod:         "POST",
			Username:             "user",
			Password:             "pass",
			SkipCertVerification: false,
		},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
			stateKey:   "component-1",
		},
	}

	rs, err := r.RemoteState()
	assert.NoError(t, err)
	assert.Equal(t, `
	data "terraform_remote_state" "component-1" {
		backend = "http"

		config = {
			address               = "https://example.com/state"
			update_method         = "POST"
			username              = "user"
			password              = "pass"
			skip_cert_verification = false
		}
	}
	`, rs)
}
