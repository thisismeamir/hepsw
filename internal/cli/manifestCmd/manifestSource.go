package manifestCmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
)

func runManifestSource(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	accessor := manifest.NewManifestAccessor(m)

	fmt.Printf("Source Information for %s@%s\n", accessor.Name(), accessor.Version())
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	source := accessor.Source()

	fmt.Printf("Type:     %s\n", source.Type)
	fmt.Printf("URL:      %s\n", source.Url)

	if source.Tag != "" {
		fmt.Printf("Tag:      %s\n", source.Tag)
	}

	if source.Checksum != "" {
		fmt.Printf("Checksum: %s\n", source.Checksum)

		// Parse checksum
		parts := strings.Split(source.Checksum, ":")
		if len(parts) == 2 {
			fmt.Printf("  Algorithm: %s\n", parts[0])
			fmt.Printf("  Hash:      %s\n", parts[1])
		}
	} else {
		fmt.Println("Checksum: (not provided)")
		fmt.Println("  Warning: Checksums are recommended for reproducibility")
	}

	return nil
}
