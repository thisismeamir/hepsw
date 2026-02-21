package manifestCmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
)

func runManifestExplain(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]
	target := args[1]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	accessor := manifest.NewManifestAccessor(m)

	fmt.Printf("Explaining: %s in %s@%s\n", target, m.Name, m.Version)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Check if it's a dependency
	if dep, found := accessor.GetDependency(target); found {
		fmt.Printf("Dependency: %s\n", dep.Name)
		fmt.Printf("Version: %s\n", dep.Version)

		if len(dep.ForOptions) > 0 {
			fmt.Printf("\nThis dependency is CONDITIONAL.\n")
			fmt.Printf("It will be included when these options are enabled:\n")
			for _, opt := range dep.ForOptions {
				fmt.Printf("  - %s\n", opt)
			}
		} else {
			fmt.Printf("\nThis dependency is ALWAYS included.\n")
		}

		if len(dep.WithOptions) > 0 {
			fmt.Printf("\nThis dependency REQUIRES these options:\n")
			for _, opt := range dep.WithOptions {
				fmt.Printf("  - %s\n", opt)
			}
		}

		if dep.IsOptional {
			fmt.Printf("\nThis dependency is OPTIONAL.\n")
		}

		return nil
	}

	// Check if it's an option
	if accessor.HasOption(target) {
		fmt.Printf("Option: %s\n\n", target)
		fmt.Println("This option affects the following dependencies:")

		affected := false
		for _, dep := range accessor.AllDependencies() {
			for _, opt := range dep.ForOptions {
				if opt == target {
					fmt.Printf("  - %s %s (enables)\n", dep.Name, dep.Version)
					affected = true
				}
			}
		}

		if !affected {
			fmt.Println("  (no dependencies are affected)")
		}

		return nil
	}

	return fmt.Errorf("target not found: %s", target)
}
