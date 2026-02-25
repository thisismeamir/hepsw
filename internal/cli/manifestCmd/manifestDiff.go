package manifestCmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
)

func runManifestDiff(cmd *cobra.Command, args []string) error {
	manifest1Path := args[0]
	manifest2Path := args[1]

	// Load both manifests
	m1, err := loader.LoadManifest(manifest1Path)
	if err != nil {
		return fmt.Errorf("failed to load first manifest: %w", err)
	}

	m2, err := loader.LoadManifest(manifest2Path)
	if err != nil {
		return fmt.Errorf("failed to load second manifest: %w", err)
	}

	// Compare manifests
	fmt.Printf("Comparing:\n  %s (%s@%s)\n  %s (%s@%s)\n\n",
		manifest1Path, m1.Name, m1.Version,
		manifest2Path, m2.Name, m2.Version)

	differences := compareManifests(m1, m2)

	if len(differences) == 0 {
		fmt.Println("No differences found.")
		return nil
	}

	fmt.Println("Differences:")
	for _, diff := range differences {
		fmt.Printf("  %s\n", diff)
	}

	return nil
}
