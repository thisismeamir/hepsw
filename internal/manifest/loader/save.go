package loader

import (
	"fmt"
	"os"

	"github.com/thisismeamir/hepsw/internal/configuration"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"gopkg.in/yaml.v3"
)

// SaveManifest saves a manifest to a file
func SaveManifest(m *manifest.Manifest, path string) error {
	data, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write manifest file: %w", err)
	}

	return nil
}

func SaveManifestFromRemote(packageIdentifier string) error {
	config, configError := configuration.GetConfiguration()
	if configError != nil {
		return fmt.Errorf("failed to load configuration: %w", configError)
	}

	thisManifest, manifestError := LoadManifestFromIndex(packageIdentifier)
	if manifestError != nil {
		return manifestError
	}

	if err := SaveManifest(thisManifest, config.Manifests+thisManifest.Name+"-"+thisManifest.Version+".yaml"); err != nil {
		return fmt.Errorf("failed to save manifest: %w", err)
	}
	return nil
}
