package manifestCmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
)

func runManifestStrip(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	// Create stripped version
	stripped := stripManifest(m, stripMinimal)

	// Generate output filename
	ext := filepath.Ext(manifestSource)
	base := strings.TrimSuffix(manifestSource, ext)
	outputPath := base + ".minimal" + ext

	// Save stripped manifest
	if err := loader.SaveManifest(stripped, outputPath); err != nil {
		return fmt.Errorf("failed to save stripped manifest: %w", err)
	}

	fmt.Printf("Stripped manifest saved to: %s\n", outputPath)

	return nil
}

func stripManifest(m *manifest.Manifest, minimal bool) *manifest.Manifest {
	stripped := &manifest.Manifest{
		Name:        m.Name,
		Version:     m.Version,
		Description: m.Description,
		Source:      m.Source,
		Recipe:      m.Recipe,
	}

	if !minimal {
		// Keep specifications if not minimal
		stripped.Specifications = m.Specifications
		stripped.Metadata = m.Metadata
	} else {
		// Minimal: only keep essential build specs
		stripped.Specifications.Build.Toolchain = m.Specifications.Build.Toolchain
	}

	return stripped
}
