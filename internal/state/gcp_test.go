package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGcpRendererBackendFull(t *testing.T) {
	r := GcpRenderer{
		state: &GcpState{
			Bucket: "bucket",
			Prefix: "prefix",
		},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
			key:        "component-1",
		},
	}

	b, err := r.Backend()
	assert.NoError(t, err)
	assert.Equal(t, `
	backend "gcs" {
	  bucket  = "bucket"
	  prefix = "prefix/test-1/component-1"
	}
	`, b)
}

func TestGcpRendererRemoteStateFull(t *testing.T) {
	r := GcpRenderer{
		state: &GcpState{
			Bucket: "bucket",
			Prefix: "prefix",
		},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
			key:        "component-1",
		},
	}

	rs, err := r.RemoteState()
	assert.NoError(t, err)
	assert.Equal(t, `
	data "terraform_remote_state" "component-1" {
	  backend = "gcp"
	
	  config = {
		  bucket  = "bucket"
		  prefix = "prefix/test-1/component-1"
	  }
	}
	`, rs)
}
