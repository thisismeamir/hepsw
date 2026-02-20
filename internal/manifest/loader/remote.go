package loader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/thisismeamir/hepsw/internal/configuration"
	"github.com/thisismeamir/hepsw/internal/index"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"gopkg.in/yaml.v3"
)

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

func LoadManifestFromIndex(packageIdentifier string) (*manifest.Manifest, error) {
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

	var thisManifest manifest.Manifest
	if err := yaml.Unmarshal(data, thisManifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest YAML: %w", err)
	}

	return &thisManifest, nil
}
