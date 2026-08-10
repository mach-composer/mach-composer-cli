package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mach-composer/mach-composer-cli/internal/cloud"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
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

		// Any other request, the commit query in particular, is unexpected.
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
	assert.True(t, cs.HasChanges())
	assert.Empty(t, cs.Changes)
}

// A component without a version to compare against yields no change set. The
// update should be skipped instead of crashing.
func TestFindSpecificUpdateWithoutChangeSet(t *testing.T) {
	branch := "main"
	organization := "acme"
	project := "ecommerce"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request %s", r.URL.Path)
	}))
	defer server.Close()

	ctx := context.Background()
	c := &config.ComponentConfig{
		Name:    "my-component",
		Version: cloud.LatestVersion,
		Branch:  branch,
	}
	cfg := &PartialConfig{
		client:   cloud.NewTestClient(server),
		filename: "main.yml",
		MachComposer: config.MachComposer{
			Cloud: config.MachComposerCloud{
				Organization: organization,
				Project:      project,
			},
		},
	}

	updates, err := findSpecificUpdate(ctx, cfg, "main.yml", c)
	assert.NoError(t, err)
	assert.NotNil(t, updates)
	assert.Empty(t, updates.updates)
	assert.False(t, updates.HasChanges())
}

// The cloud lookup runs on the same worker pool as the git lookup. Run with
// -race to cover the per-component logging context.
func TestFindUpdatesCloudConcurrent(t *testing.T) {
	branch := "main"
	organization := "acme"
	project := "ecommerce"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/latest") {
			t.Errorf("unexpected request %s", r.URL.Path)
			return
		}

		component := path.Base(path.Dir(r.URL.Path))
		id, _ := uuid.NewUUID()
		b, _ := json.Marshal(mccsdk.ComponentVersion{
			Id:        id.String(),
			CreatedAt: time.Now(),
			Component: component,
			Version:   component + "-new",
			Branch:    &branch,
		})
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}))
	defer server.Close()

	components := make([]config.ComponentConfig, 25)
	for i := range components {
		components[i] = config.ComponentConfig{
			Name:    fmt.Sprintf("component-%d", i),
			Version: fmt.Sprintf("component-%d-old", i),
		}
	}

	cfg := &PartialConfig{
		client:     cloud.NewTestClient(server),
		filename:   "main.yml",
		Components: components,
		MachComposer: config.MachComposer{
			Cloud: config.MachComposerCloud{
				Organization: organization,
				Project:      project,
			},
		},
	}

	updates, err := findUpdates(context.Background(), cfg, "main.yml")
	assert.NoError(t, err)
	assert.Len(t, updates.updates, len(components))

	// Each worker writes the default branch back into its own component, so the
	// defaults are visible in the shared config afterwards.
	for i := range cfg.Components {
		assert.Equal(t, branch, cfg.Components[i].Branch)
	}
}
