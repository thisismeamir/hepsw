package manifestCmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
	"gopkg.in/yaml.v3"
)

func runManifestFormat(cmd *cobra.Command, args []string) error {
	manifestPath := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	// Format and save
	data, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	// Write back to file
	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write formatted manifest: %w", err)
	}

	fmt.Printf("Formatted manifest: %s\n", manifestPath)
	return nil
}
