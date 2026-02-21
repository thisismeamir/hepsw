package manifestCmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
)

func runManifestFetch(cmd *cobra.Command, args []string) error {
	reference := args[0]

	fmt.Printf("Fetching manifest: %s\n", reference)

	// Load from index
	m, err := loader.LoadManifestFromIndex(reference)
	if err != nil {
		return fmt.Errorf("failed to fetch manifest: %w", err)
	}

	// Determine output filename
	filename := fmt.Sprintf("%s-%s.yaml", m.Name, m.Version)
	outputPath := filepath.Join(fetchDestination, filename)

	// Save manifest
	if err := loader.SaveManifest(m, outputPath); err != nil {
		return fmt.Errorf("failed to save manifest: %w", err)
	}

	fmt.Printf("Downloaded to: %s\n", outputPath)
	fmt.Printf("Package: %s@%s\n", m.Name, m.Version)

	return nil
}
