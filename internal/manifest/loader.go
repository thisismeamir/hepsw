package manifest

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/thisismeamir/hepsw/internal/utils"
	"gopkg.in/yaml.v3"
)

// LoadManifest loads a manifest from a file path or HepSW index reference
func LoadManifest(source string) (*Manifest, error) {
	// Check if it's a file path
	if utils.IsFilePath(source) {
		return LoadManifestFromFile(source)
	}

	// Otherwise treat it as an index reference
	return LoadManifestFromIndex(source)
}

// LoadManifestFromFile loads a manifest from a local file
func LoadManifestFromFile(path string) (*Manifest, error) {
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
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse manifest YAML: %w", err)
	}

	return &m, nil
}

// LoadManifestFromIndex loads a manifest from the HepSW Package Index Repository
func LoadManifestFromIndex(reference string) (*Manifest, error) {
	// Parse reference: package@version or just package
	parts := strings.Split(reference, "@")
	packageName := parts[0]
	version := "latest"
	if len(parts) > 1 {
		version = parts[1]
	}

	var m Manifest
	// Get index URL from config or environment
	indexURL := getIndexURL()

	// Construct manifest URL
	manifestURL := filepath.Join(indexURL, packageName, version+".yaml")

	if utils.IsFilePath(manifestURL) {
		// Open the file of the manifest (This is a faster route if the index repository is locally available.
		data, err := os.ReadFile(manifestURL)
		if err != nil {
			return nil, fmt.Errorf("failed to read manifest file: %w", err)
		}
		if err := yaml.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("failed to parse manifest YAML from index: %w", err)
		}
		return &m, nil
	} else {
		// Fetch manifest
		resp, err := http.Get(manifestURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch manifest from index: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("manifest not found in index: %s (status: %d)", reference, resp.StatusCode)
		}

		// Read response
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read manifest from index: %w", err)
		}

		// Parse YAML

		if err := yaml.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("failed to parse manifest YAML from index: %w", err)
		}

		return &m, nil

	}
}

// LoadManifestFromReader loads a manifest from an io.Reader
func LoadManifestFromReader(r io.Reader) (*Manifest, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse manifest YAML: %w", err)
	}

	return &m, nil
}

// SaveManifest saves a manifest to a file
func SaveManifest(m *Manifest, path string) error {
	data, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write manifest file: %w", err)
	}

	return nil
}

// getIndexURL returns the HepSW Package Index Repository URL
func getIndexURL() string {
	// Try environment variable first
	if url := os.Getenv("HEPSW_INDEX_URL"); url != "" {
		return url
	}

	// Try config file
	// TODO: Implement config file reading

	// Default to official repository (To be used later
	//return "https://index.hepsw.org"
	return filepath.Join(os.UserHomeDir(), "index.yaml")
}
