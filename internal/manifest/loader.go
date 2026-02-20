package manifest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/thisismeamir/hepsw/internal/configuration"
	"github.com/thisismeamir/hepsw/internal/index"
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

func getManifestURL(packageName string, version string) (*string, error) {
	config, configurationError := configuration.GetConfiguration()
	if configurationError != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", configurationError)
	}
	newIndex, indexGenerationError := index.New(&config.IndexConfig)
	if indexGenerationError != nil {
		return nil, fmt.Errorf("failed to load index: %w", indexGenerationError)
	}

	manifestIndexEntry, manifestIndexEntryError := newIndex.GetPackage(context.Background(), packageName)
	if manifestIndexEntryError != nil {
		if strings.Contains(manifestIndexEntryError.Error(), "not found") {
			return nil, fmt.Errorf("package not found in index: %s", packageName)
		}
		return nil, fmt.Errorf("failed to get package from index: %w", manifestIndexEntryError)
	}

	if version == "latest" {
		manifest, fetchingURLError := newIndex.GetLatestVersion(context.Background(), manifestIndexEntry.Name)
		if fetchingURLError != nil {
			return nil, fmt.Errorf("failed to fetch latest version URL: %w", fetchingURLError)
		}
		return &manifest.ManifestURL, nil
	} else {
		manifest, fetchingURLError := newIndex.GetVersion(context.Background(), manifestIndexEntry.Name, version)
		if fetchingURLError != nil {
			return nil, fmt.Errorf("failed to fetch version URL: %w", fetchingURLError)
		}
		return &manifest.ManifestURL, nil
	}

}
func LoadManifestFromIndex(packageIdentifier string) (*Manifest, error) {
	// Separate by @ to get package name and version
	parts := strings.SplitN(packageIdentifier, "@", 2)
	var packageName string
	var version string
	if len(parts) == 1 {
		// assuming the latest package:
		packageName = parts[0]
		version = "latest"
	} else if len(parts) == 2 {
		packageName = parts[0]
		version = parts[1]
	} else {
		return nil, fmt.Errorf("invalid package identifier format: %s", packageIdentifier)
	}

	//get manifest URL from index
	manifestURL, urlFetchingError := getManifestURL(packageName, version)
	if urlFetchingError != nil {
		return nil, urlFetchingError
	}

	// Checking the validity of the URL (it must be YAML).
	if !strings.HasSuffix(strings.ToLower(*manifestURL), ".yaml") &&
		!strings.HasSuffix(strings.ToLower(*manifestURL), ".yml") {
		return nil, fmt.Errorf("URL does not point to a YAML file: %s", manifestURL)
	}

	resp, err := http.Get(*manifestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch manifest from index: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch manifest: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest from index: %w", err)
	}

	manifestPath, err := saveManifest(packageIdentifier, data)
	if err != nil {
		return nil, err
	}

	return ReadManifest(manifestPath)
}

func saveManifest(packageIdentifier string, data []byte) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine user home directory: %w", err)
	}

	manifestsDir := filepath.Join(homeDir, ".hepsw", "manifests")
	if err := os.MkdirAll(manifestsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create manifests directory: %w", err)
	}

	manifestPath := filepath.Join(manifestsDir, packageIdentifier+".yaml")
	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to save manifest to disk: %w", err)
	}

	return manifestPath, nil
}

func ReadManifest(manifestPath string) (*Manifest, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest from disk: %w", err)
	}

	manifest := &Manifest{}
	if err := yaml.Unmarshal(data, manifest); err != nil {
		return nil, fmt.Errorf("failed to unmarshal manifest: %w", err)
	}

	return manifest, nil
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
