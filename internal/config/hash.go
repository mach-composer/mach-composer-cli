package config

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/mach-composer/mach-composer-cli/internal/variables"
)

func ComputeHash(cfg *MachConfig) (string, error) {
	hashConfig := struct {
		MachComposer MachComposer         `json:"mach_composer"`
		Global       GlobalConfig         `json:"global"`
		Sites        []SiteConfig         `json:"sites"`
		Components   []ComponentConfig    `json:"components"`
		Filename     string               `json:"filename"`
		Variables    *variables.Variables `json:"variables"`
	}{
		MachComposer: cfg.MachComposer,
		Global:       cfg.Global,
		Sites:        cfg.Sites,
		Components:   cfg.Components,
		Filename:     cfg.Filename,
		Variables:    cfg.Variables,
	}
	data, err := json.Marshal(hashConfig)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil)), nil
}
