package hash

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"os"
)

const defaultHashFile = ".mach-composer/hashes.json"

type Handler interface {
	Store(ctx context.Context, n graph.Node) error
	Fetch(ctx context.Context, n graph.Node) (string, error)
}

func Factory(_ *config.MachConfig) Handler {
	hashFile := os.Getenv("MC_HASH_FILE")
	if hashFile == "" {
		hashFile = defaultHashFile
	}

	return NewJsonFileHandler(hashFile)
}
