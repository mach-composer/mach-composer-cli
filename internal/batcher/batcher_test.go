package batcher

import (
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReturnsErrorWhenUnknownBatchType(t *testing.T) {
	cfg := &config.MachConfig{
		MachComposer: config.MachComposer{
			Batcher: config.Batcher{
				Type: "unknown",
			},
		},
	}

	_, err := Factory(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown batch type unknown")
}

func TestReturnsSimpleBatchFuncWhenTypeIsEmpty(t *testing.T) {
	cfg := &config.MachConfig{
		MachComposer: config.MachComposer{
			Batcher: config.Batcher{
				Type: "",
			},
		},
	}

	batchFunc, err := Factory(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, batchFunc)
}

func TestReturnsSimpleBatchFuncWhenTypeIsSimple(t *testing.T) {
	cfg := &config.MachConfig{
		MachComposer: config.MachComposer{
			Batcher: config.Batcher{
				Type: "simple",
			},
		},
	}

	batchFunc, err := Factory(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, batchFunc)
}

func TestReturnsSiteBatchFunc(t *testing.T) {
	cfg := &config.MachConfig{
		MachComposer: config.MachComposer{
			Batcher: config.Batcher{
				Type: "site",
			},
		},
		Sites: config.SiteConfigs{
			{
				Identifier: "site-1",
			},
			{
				Identifier: "site-2",
			},
		},
	}

	batchFunc, err := Factory(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, batchFunc)
}

func TestDetermineSiteOrderReturnsIdentifiersWhenNoSiteOrderProvided(t *testing.T) {
	cfg := &config.MachConfig{
		MachComposer: config.MachComposer{
			Batcher: config.Batcher{},
		},
		Sites: config.SiteConfigs{
			{Identifier: "site-1"},
			{Identifier: "site-2"},
		},
	}
	order, err := DetermineSiteOrder(cfg)
	assert.NoError(t, err)
	assert.Equal(t, []string{"site-1", "site-2"}, order)
}

func TestDetermineSiteOrderReturnsSiteOrderWhenProvided(t *testing.T) {
	cfg := &config.MachConfig{
		MachComposer: config.MachComposer{
			Batcher: config.Batcher{
				SiteOrder: []string{"site-2", "site-1"},
			},
		},
		Sites: config.SiteConfigs{
			{Identifier: "site-1"},
			{Identifier: "site-2"},
		},
	}
	order, err := DetermineSiteOrder(cfg)
	assert.NoError(t, err)
	assert.Equal(t, []string{"site-2", "site-1"}, order)
}

func TestDetermineSiteOrderReturnsErrorWhenSiteOrderLengthMismatch(t *testing.T) {
	cfg := &config.MachConfig{
		MachComposer: config.MachComposer{
			Batcher: config.Batcher{
				SiteOrder: []string{"site-1", "site-2"},
			},
		},
		Sites: config.SiteConfigs{
			{Identifier: "site-1"},
		},
	}
	order, err := DetermineSiteOrder(cfg)
	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Contains(t, err.Error(), "site order length 2 does not match identifiers length 1")
}

func TestDetermineSiteOrderReturnsErrorWhenSiteOrderContainsUnknownIdentifier(t *testing.T) {
	cfg := &config.MachConfig{
		MachComposer: config.MachComposer{
			Batcher: config.Batcher{
				SiteOrder: []string{"site-1", "unknown-site"},
			},
		},
		Sites: config.SiteConfigs{
			{Identifier: "site-1"},
			{Identifier: "site-2"},
		},
	}
	order, err := DetermineSiteOrder(cfg)
	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Contains(t, err.Error(), "site order contains siteIdentifier unknown-site that is not in the identifiers list")
}

func TestDetermineSiteOrderReturnsEmptyWhenNoSites(t *testing.T) {
	cfg := &config.MachConfig{
		MachComposer: config.MachComposer{
			Batcher: config.Batcher{},
		},
		Sites: config.SiteConfigs{},
	}
	order, err := DetermineSiteOrder(cfg)
	assert.NoError(t, err)
	assert.Empty(t, order)
}
