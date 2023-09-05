package state

import "fmt"

type Repository struct {
	renderers map[string]Renderer
}

func NewRepository() *Repository {
	return &Repository{
		renderers: make(map[string]Renderer),
	}
}

func (r *Repository) Add(key string, renderer Renderer) error {
	if key == "" {
		return fmt.Errorf("renderer key cannot be empty")
	}

	r.renderers[key] = renderer
	return nil
}

func (r *Repository) Has(key string) bool {
	_, ok := r.renderers[key]
	return ok
}

func (r *Repository) Get(key string) Renderer {
	renderer, ok := r.renderers[key]
	if !ok {
		return nil
	}
	return renderer
}
