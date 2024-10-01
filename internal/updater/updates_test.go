package updater

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mach-composer/mach-composer-cli/internal/cloud"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetLastVersionCloudLatestVersion(t *testing.T) {
	branch := "main"
	organization := "acme"
	project := "ecommerce"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request %s", r.URL.Path)
	}))
	defer server.Close()

	ctx := context.Background()
	c := &config.ComponentConfig{
		Version: cloud.LatestVersion,
		Branch:  branch,
	}
	cfg := &PartialConfig{
		client: cloud.NewTestClient(server),
		MachComposer: config.MachComposer{
			Cloud: config.MachComposerCloud{
				Organization: organization,
				Project:      project,
			},
		},
	}

	cs, err := getLastVersionCloud(ctx, cfg, c, branch)
	assert.NoError(t, err)
	assert.Nil(t, cs)
}

func TestGetLastVersionCloudVersionNotApplicable(t *testing.T) {
	branch := "main"
	organization := "acme"
	project := "ecommerce"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request %s", r.URL.Path)
	}))
	defer server.Close()

	ctx := context.Background()
	c := &config.ComponentConfig{
		Version: cloud.VersionNotApplicable,
		Branch:  branch,
	}
	cfg := &PartialConfig{
		client: cloud.NewTestClient(server),
		MachComposer: config.MachComposer{
			Cloud: config.MachComposerCloud{
				Organization: organization,
				Project:      project,
			},
		},
	}

	cs, err := getLastVersionCloud(ctx, cfg, c, branch)
	assert.NoError(t, err)
	assert.Nil(t, cs)
}

func TestGetLastVersionCloudOK(t *testing.T) {
	branch := "main"
	organization := "acme"
	project := "ecommerce"
	component := "component"
	newVersion := "0.0.2"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/organizations/acme/projects/ecommerce/components//latest" {
			id, _ := uuid.NewUUID()
			b, _ := json.Marshal(mccsdk.ComponentVersion{
				Id:        id.String(),
				CreatedAt: time.Now(),
				Component: component,
				Version:   newVersion,
				Branch:    &branch,
			})
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(b)
			return
		}

		if r.URL.Path == "/organizations/acme/projects/ecommerce/components//commits" {
			b, _ := json.Marshal(mccsdk.CommitDataPaginator{
				Count: utils.Ref(int32(1)),
				Total: utils.Ref(int64(1)),
				Results: []mccsdk.CommitData{
					{
						Commit:    "test",
						Parents:   nil,
						Subject:   "test",
						Author:    mccsdk.CommitDataAuthor{},
						Committer: mccsdk.CommitDataAuthor{},
						//TODO: should these fields be available?
						//Version:   &newVersion,
						//Branch:    &branch,
					},
				},
			})
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(b)
			return
		}

		t.Errorf("unexpected request %s", r.URL.Path)
	}))
	defer server.Close()

	ctx := context.Background()
	c := &config.ComponentConfig{
		Version: "0.0.1",
		Branch:  branch,
	}
	cfg := &PartialConfig{
		client: cloud.NewTestClient(server),
		MachComposer: config.MachComposer{
			Cloud: config.MachComposerCloud{
				Organization: organization,
				Project:      project,
			},
		},
	}

	cs, err := getLastVersionCloud(ctx, cfg, c, branch)
	assert.NoError(t, err)
	assert.NotNil(t, cs)
	assert.Equal(t, newVersion, cs.LastVersion)
	assert.Equal(t, 1, len(cs.Changes))
}
