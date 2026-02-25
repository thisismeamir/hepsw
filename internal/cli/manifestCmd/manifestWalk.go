package manifestCmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
)

func runManifestWalk(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	// Walk manifest
	result, err := manifest.WalkManifest(m, walkOptions, walkVariables)
	if err != nil {
		return fmt.Errorf("failed to walk manifest: %w", err)
	}

	// Print result
	output := manifest.PrintWalkResult(result)
	fmt.Print(output)

	return nil
}
