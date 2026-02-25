package loader

import (
	"fmt"
	"io"

	"github.com/thisismeamir/hepsw/internal/manifest"
	"github.com/thisismeamir/hepsw/internal/utils"
	"gopkg.in/yaml.v3"
)

// LoadManifest loads a manifest from a file path or HepSW index reference
func LoadManifest(source string) (*manifest.Manifest, error) {
	// Check if it's a file path
	if utils.IsFilePath(source) {
		return LoadManifestFromFile(source)
	}

	// Otherwise treat it as an index reference
	return LoadManifestFromIndex(source)
}

// LoadManifestFromReader loads a manifest from an io.Reader
func LoadManifestFromReader(r io.Reader) (*manifest.Manifest, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var m manifest.Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse manifest YAML: %w", err)
	}

	return &m, nil
}
