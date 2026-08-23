// Package manifest provides a programmatic way to build and validate
// Flow Launcher plugin.json manifests.
package manifest

import (
	"encoding/json"
	"fmt"
)

// Config defines the fields for a Flow Launcher plugin manifest.
type Config struct {
	ID            string
	ActionKeyword string
	Name          string
	Description   string
	Author        string
	Version       string
	Language      string
	Website       string
	IcoPath       string
	ExeFileName   string
}

// Manifest represents a Flow Launcher plugin.json structure.
type Manifest struct {
	Schema          string `json:"$schema,omitempty"`
	ID              string `json:"ID"`
	ActionKeyword   string `json:"ActionKeyword"`
	Name            string `json:"Name"`
	Description     string `json:"Description"`
	Author          string `json:"Author"`
	Version         string `json:"Version"`
	Language        string `json:"Language"`
	Website         string `json:"Website,omitempty"`
	IcoPath         string `json:"IcoPath,omitempty"`
	ExecuteFileName string `json:"ExecuteFileName"`
}

// New creates a Manifest with the given config and sensible defaults.
func New(cfg Config) *Manifest {
	m := &Manifest{
		Schema:          "https://www.flowlauncher.com/schemas/plugin.schema.json",
		ID:              cfg.ID,
		ActionKeyword:   cfg.ActionKeyword,
		Name:            cfg.Name,
		Description:     cfg.Description,
		Author:          cfg.Author,
		Version:         cfg.Version,
		Language:        cfg.Language,
		Website:         cfg.Website,
		IcoPath:         cfg.IcoPath,
		ExecuteFileName: cfg.ExeFileName,
	}
	if m.Language == "" {
		m.Language = "executable"
	}
	if m.Version == "" {
		m.Version = "1.0.0"
	}
	return m
}

// Validate checks that required fields are set.
func (m *Manifest) Validate() error {
	if m.ID == "" {
		return fmt.Errorf("manifest: ID is required")
	}
	if m.ActionKeyword == "" {
		return fmt.Errorf("manifest: ActionKeyword is required")
	}
	if m.Name == "" {
		return fmt.Errorf("manifest: Name is required")
	}
	if m.Description == "" {
		return fmt.Errorf("manifest: Description is required")
	}
	if m.Author == "" {
		return fmt.Errorf("manifest: Author is required")
	}
	if m.Version == "" {
		return fmt.Errorf("manifest: Version is required")
	}
	if m.Language == "" {
		return fmt.Errorf("manifest: Language is required")
	}
	if m.ExecuteFileName == "" {
		return fmt.Errorf("manifest: ExecuteFileName is required")
	}
	return nil
}

// Marshal returns the pretty-printed JSON of the manifest.
func (m *Manifest) Marshal() ([]byte, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return json.MarshalIndent(m, "", "  ")
}

// Unmarshal parses a manifest from JSON bytes.
func Unmarshal(data []byte) (*Manifest, error) {
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("manifest: invalid JSON: %w", err)
	}
	return &m, nil
}