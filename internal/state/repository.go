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

func (r *Repository) Add(renderer Renderer) error {
	if renderer.Identifier() == "" {
		return fmt.Errorf("renderer Identifier cannot be empty")
	}

	r.renderers[renderer.Identifier()] = renderer
	return nil
}

func (r *Repository) Key(identifier string) (string, bool) {
	rr, ok := r.renderers[identifier]
	if ok {
		return rr.Key(), true
	}

	for alias, k := range r.aliases {
		if alias == identifier {
			return k, true
		}
	}

	return "", false
}

func (r *Repository) Has(identifier string) bool {
	_, ok := r.renderers[identifier]
	return ok
}

func (r *Repository) Get(identifier string) (Renderer, bool) {
	return r.renderers[identifier], r.Has(identifier)
}

func (r *Repository) Alias(identifier string, alias string) {
	r.aliases[alias] = identifier
}
