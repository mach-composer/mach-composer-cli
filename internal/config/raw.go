package config

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/plugins"
	"gopkg.in/yaml.v3"
)

type rawConfig struct {
	MachComposer MachComposer `yaml:"mach_composer"`
	Global       yaml.Node    `yaml:"global"`
	Sites        yaml.Node    `yaml:"sites"`
	Components   yaml.Node    `yaml:"components"`

	document  *yaml.Node                `yaml:"-"`
	filename  string                    `yaml:"-"`
	plugins   *plugins.PluginRepository `yaml:"-"`
	variables *Variables                `yaml:"-"`
}

func (r *rawConfig) validate() error {
	if r.MachComposer.Version == "" {
		return fmt.Errorf("no version")
	}

	if r.filename == "" {
		return fmt.Errorf("filename must be set")
	}
	if r.variables == nil {
		return fmt.Errorf("variables cannot be nil")
	}
	if r.plugins == nil {
		return fmt.Errorf("plugins cannot be nil")
	}

	return nil
}

func (r *rawConfig) computeHash() (string, error) {
	hashConfig := struct {
		MachComposer MachComposer `json:"mach_composer"`
		Global       yaml.Node    `json:"global"`
		Sites        yaml.Node    `json:"sites"`
		Components   yaml.Node    `json:"components"`
		Filename     string       `json:"filename"`
		Variables    *Variables   `json:"variables"`
	}{
		MachComposer: r.MachComposer,
		Global:       r.Global,
		Sites:        r.Sites,
		Components:   r.Components,
		Filename:     r.filename,
		Variables:    r.variables,
	}
	data, err := json.Marshal(hashConfig)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil)), nil
}

func newRawConfig(filename string, document *yaml.Node) (*rawConfig, error) {
	r := &rawConfig{
		filename:  filename,
		variables: NewVariables(),
		document:  document,
	}
	if err := document.Decode(r); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}
	return r, nil
}
