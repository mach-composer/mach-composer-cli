package state

import (
	"fmt"
)

type Repository struct {
	aliases   map[string]string
	renderers map[string]Renderer
}

func NewRepository() *Repository {
	return &Repository{
		aliases:   map[string]string{},
		renderers: make(map[string]Renderer),
	}
}

func (r *Repository) Add(key string, renderer Renderer) error {
	if key == "" {
		return fmt.Errorf("renderer Key cannot be empty")
	}

	r.renderers[key] = renderer
	return nil
}

func (r *Repository) Key(key string) (string, bool) {
	_, ok := r.renderers[key]
	if ok {
		return key, true
	}

	for alias, k := range r.aliases {
		if alias == key {
			return k, true
		}
	}

	return "", false
}

func (r *Repository) Has(key string) bool {
	_, ok := r.Key(key)
	return ok
}

func (r *Repository) Get(key string) Renderer {
	k, ok := r.Key(key)
	if !ok {
		return nil
	}
	return r.renderers[k]
}

func (r *Repository) Alias(key string, alias string) {
	r.aliases[alias] = key
}
