package manifestCmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
	"gopkg.in/yaml.v3"
)

func runManifestFlatten(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	// Create flattened version
	flattened := flattenManifest(m, flattenOptions)

	// Output
	data, err := yaml.Marshal(flattened)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	if flattenOutput != "" {
		if err := os.WriteFile(flattenOutput, data, 0644); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		fmt.Printf("Flattened manifest written to: %s\n", flattenOutput)
	} else {
		fmt.Print(string(data))
	}

	return nil
}

func flattenManifest(m *manifest.Manifest, options []string) *manifest.Manifest {
	// Create a copy
	flattened := *m

	// Filter dependencies based on options
	accessor := manifest.NewManifestAccessor(m)
	flattened.Specifications.Build.Dependencies = accessor.GetDependenciesForOptions(options)

	// Remove conditional steps from recipe
	// (This is a simplified version - full implementation would evaluate all conditionals)

	return &flattened
}
