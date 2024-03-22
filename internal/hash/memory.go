package hash

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
)

type Entry struct {
	Identifier string
	Hash       string
}

type MemoryMap struct {
	InternalMap map[string]string
}

func NewMemoryMapHandler(entries ...Entry) Handler {
	h := &MemoryMap{
		InternalMap: make(map[string]string),
	}

	for _, e := range entries {
		h.InternalMap[e.Identifier] = e.Hash
	}

	return h
}

func (h *MemoryMap) Fetch(_ context.Context, n graph.Node) (string, error) {
	return h.InternalMap[n.Identifier()], nil
}
func (h *MemoryMap) Store(_ context.Context, n graph.Node) error {
	var err error
	h.InternalMap[n.Identifier()], err = n.Hash()

	return err
}
