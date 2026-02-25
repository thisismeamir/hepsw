package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/thisismeamir/hepsw/internal/manifest"
	"gopkg.in/yaml.v3"
)

// LoadManifestFromFile loads a manifest from a local file
func LoadManifestFromFile(path string) (*manifest.Manifest, error) {
	// Resolve absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("manifest file not found: %s", absPath)
	}

	// Read file
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest file: %w", err)
	}

	// Parse YAML
	var m manifest.Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse manifest YAML: %w", err)
	}

	return &m, nil
}
