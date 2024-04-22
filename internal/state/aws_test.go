package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAwsRenderer(t *testing.T) {
	r := AwsRenderer{
		state: &AwsState{
			Bucket:    "bucket",
			KeyPrefix: "key",
			Region:    "region",
			RoleARN:   "role",
			LockTable: "lock",
			Encrypt:   true,
		},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
			stateKey:   "component-1",
		},
	}

	b, err := r.Backend()
	assert.NoError(t, err)
	assert.Equal(t, `
	backend "s3" {
	  bucket         = "bucket"
	  key            = "key/test-1/component-1"
	  region         = "region"
	  
	  role_arn       = "role"
	  
	  
	  dynamodb_table = "lock"
	  
	  encrypt        = true
	}
	`, b)
}

func TestAwsRendererRemoteStateFull(t *testing.T) {
	r := AwsRenderer{
		state: &AwsState{
			Bucket:    "bucket",
			KeyPrefix: "key",
			Region:    "region",
			RoleARN:   "role",
			LockTable: "lock",
			Encrypt:   true,
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
	  backend = "s3"

	  config = {
		  bucket         = "bucket"
		  key            = "key/test-1/component-1"
		  region         = "region"
		  
		  role_arn       = "role"
		  
		  
		  dynamodb_table = "lock"
		  
		  encrypt        = true
	  }
	}
	`, rs)
}
