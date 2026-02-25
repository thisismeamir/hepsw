package manifestCmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
)

func runManifestOptions(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	accessor := manifest.NewManifestAccessor(m)

	fmt.Printf("Options for %s@%s\n", accessor.Name(), accessor.Version())
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	options := accessor.Options()
	if len(options) == 0 {
		fmt.Println("No build options defined.")
		return nil
	}

	// Show each option and its effects
	for _, opt := range options {
		fmt.Printf("Option: %s\n", opt)

		// Find dependencies enabled by this option
		enabledDeps := make([]string, 0)
		for _, dep := range accessor.AllDependencies() {
			for _, reqOpt := range dep.ForOptions {
				if reqOpt == opt {
					enabledDeps = append(enabledDeps, fmt.Sprintf("%s %s", dep.Name, dep.Version))
					break
				}
			}
		}

		if len(enabledDeps) > 0 {
			fmt.Println("  Enables dependencies:")
			for _, dep := range enabledDeps {
				fmt.Printf("    - %s\n", dep)
			}
		} else {
			fmt.Println("  No dependencies affected")
		}
		fmt.Println()
	}

	return nil
}
