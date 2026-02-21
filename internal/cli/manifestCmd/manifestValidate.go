package manifestCmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
)

func runManifestValidate(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	fmt.Printf("Validating manifest: %s@%s\n", m.Name, m.Version)
	fmt.Println()

	// Validate
	result := manifest.ValidateManifest(m)

	// Print results
	if len(result.Errors) > 0 {
		fmt.Println("ERRORS:")
		for _, err := range result.Errors {
			fmt.Printf("  ✗ %s\n", err)
		}
		fmt.Println()
	}

	if len(result.Warnings) > 0 {
		fmt.Println("WARNINGS:")
		for _, warn := range result.Warnings {
			fmt.Printf("  ⚠ %s\n", warn)
		}
		fmt.Println()
	}

	if validateVerbose && len(result.Info) > 0 {
		fmt.Println("INFO:")
		for _, info := range result.Info {
			fmt.Printf("  ℹ %s\n", info)
		}
		fmt.Println()
	}

	// Summary
	if result.Valid {
		fmt.Println("✓ Manifest is valid")
		return nil
	} else {
		fmt.Printf("✗ Manifest validation failed with %d error(s)\n", len(result.Errors))
		return fmt.Errorf("validation failed")
	}
}
