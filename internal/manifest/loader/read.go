package loader

import (
	"fmt"
	"os"

	"github.com/thisismeamir/hepsw/internal/manifest"
	"gopkg.in/yaml.v3"
)

func ReadManifest(manifestPath string) (*manifest.Manifest, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest from disk: %w", err)
	}

	manifest := &manifest.Manifest{}
	if err := yaml.Unmarshal(data, manifest); err != nil {
		return nil, fmt.Errorf("failed to unmarshal manifest: %w", err)
	}

	return manifest, nil
}
