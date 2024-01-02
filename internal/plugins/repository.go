package plugins

import (
	"context"
	"fmt"
	"github.com/elliotchance/pie/v2"
	"strings"

	schemav2 "github.com/mach-composer/mach-composer-plugin-sdk/v2/schema"
	"github.com/rs/zerolog/log"
)

type PluginNotFoundError struct {
	name string
}

func (e *PluginNotFoundError) Error() string {
	return fmt.Sprintf("plugin %s not found", e.name)
}

type PluginError struct {
	msg    string
	plugin string
}

func (e *PluginError) Error() string {
	return e.msg
}

type PluginRepository struct {
	handlers map[string]*PluginHandler
	schemas  map[string]schemav2.ValidationSchema
}

func NewPluginRepository() *PluginRepository {
	return &PluginRepository{
		handlers: make(map[string]*PluginHandler),
		schemas:  make(map[string]schemav2.ValidationSchema),
	}
}

// Close kills all the running handlers
func (p *PluginRepository) Close() {
	for _, rp := range p.handlers {
		rp.Close()
	}
}

// All returns all the plugin names in the repository, ordered by the plugin name
func (p *PluginRepository) All() []*PluginHandler {
	result := make([]*PluginHandler, len(p.handlers))
	for i, key := range pie.Sort(pie.Keys(p.handlers)) {
		result[i] = p.handlers[key]
	}
	return result
}

func (p *PluginRepository) Names(names ...string) []PluginHandler {
	var result []PluginHandler

	for _, key := range pie.Sort(pie.Keys(p.handlers)) {
		if !pie.Contains(names, key) {
			continue
		}

		result = append(result, *p.handlers[key])
	}

	return result
}

// Add an existing plugin to the repository. Only used for testing purposes
func (p *PluginRepository) Add(name string, machComposerPlugin schemav2.MachComposerPlugin) error {
	handlers := &PluginHandler{
		MachComposerPlugin: machComposerPlugin,
		Name:               name,
		isRunning:          true,
	}

	p.handlers[name] = handlers
	return nil
}

// Get returns the plugin from the repository
func (p *PluginRepository) Get(name string) (*PluginHandler, error) {
	if name == "" {
		panic("handler name is empty") // this should never happen
	}
	handler, ok := p.handlers[name]
	if !ok {
		return nil, &PluginNotFoundError{name: name}
	}
	return handler, nil
}

// LoadPlugin will load the plugin with the given name and Start it
func (p *PluginRepository) LoadPlugin(ctx context.Context, name string, config PluginConfig) error {
	if _, ok := p.handlers[name]; ok {
		return fmt.Errorf("handler %s is already loaded", name)
	}

	handler := &PluginHandler{
		Name:   name,
		Config: config,
	}
	p.handlers[name] = handler

	if err := handler.Start(ctx); err != nil {
		return err
	}
	return nil
}

func (p *PluginRepository) loadSchema(name string) (*schemav2.ValidationSchema, error) {
	plugin, err := p.Get(name)
	if err != nil {
		return nil, err
	}

	// Load validation schema
	schema, err := plugin.GetValidationSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin schema %s: %w", name, err)
	}
	if schema != nil {
		p.schemas[name] = *schema
	}
	return schema, nil
}

func (p *PluginRepository) GetSchema(name string) (*schemav2.ValidationSchema, error) {
	if name == "" {
		panic("plugin name is empty") // this should never happen
	}

	// this should not happen in a regular use-case
	if schema, ok := p.schemas[name]; ok {
		return &schema, nil
	}

	return p.loadSchema(name)
}

func (p *PluginRepository) handleError(plugin string, err error) error {
	if err == nil {
		return nil
	}

	log.Error().Err(err).Stack().Msgf("plugin %s returned an error", plugin)

	if strings.Contains(err.Error(), "reading body EOF") {
		return &PluginError{
			plugin: plugin,
			msg:    fmt.Sprintf("the %s plugin crashed. This is a bug in the plugin.", plugin),
		}
	}
	return &PluginError{
		plugin: plugin,
		msg:    err.Error(),
	}
}
