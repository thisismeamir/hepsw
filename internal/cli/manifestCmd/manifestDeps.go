package manifestCmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
)

func runManifestDeps(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	accessor := manifest.NewManifestAccessor(m)

	fmt.Printf("Dependencies for %s@%s\n", accessor.Name(), accessor.Version())
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Build dependencies
	buildDeps := accessor.BuildDependencies()
	if len(buildDeps) > 0 {
		fmt.Println("Build Dependencies:")
		for _, dep := range buildDeps {
			fmt.Printf("  ├─ %s %s", dep.Name, dep.Version)
			if depsShowOptions && len(dep.ForOptions) > 0 {
				fmt.Printf(" [for: %s]", strings.Join(dep.ForOptions, ", "))
			}
			if dep.IsOptional {
				fmt.Print(" (optional)")
			}
			fmt.Println()
		}
		fmt.Println()
	}

	// Runtime dependencies
	runtimeDeps := accessor.RuntimeDependencies()
	if len(runtimeDeps) > 0 {
		fmt.Println("Runtime Dependencies:")
		for _, dep := range runtimeDeps {
			fmt.Printf("  ├─ %s %s", dep.Name, dep.Version)
			if depsShowOptions && len(dep.ForOptions) > 0 {
				fmt.Printf(" [for: %s]", strings.Join(dep.ForOptions, ", "))
			}
			fmt.Println()
		}
		fmt.Println()
	}

	if len(buildDeps) == 0 && len(runtimeDeps) == 0 {
		fmt.Println("No dependencies found.")
	}

	return nil
}
