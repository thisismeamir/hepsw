package manifestCmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
)

func runManifestRecipe(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	accessor := manifest.NewManifestAccessor(m)

	fmt.Printf("Recipe for %s@%s\n", accessor.Name(), accessor.Version())
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Print each phase
	phases := []struct {
		name  string
		steps []manifest.RecipeStep
	}{
		{"Configuration", accessor.ConfigurationSteps()},
		{"Build", accessor.BuildSteps()},
		{"Install", accessor.InstallSteps()},
		{"Use", accessor.UseSteps()},
	}

	for _, phase := range phases {
		if len(phase.steps) == 0 {
			continue
		}

		fmt.Printf("=== %s Phase ===\n\n", phase.name)

		for i, step := range phase.steps {
			fmt.Printf("%d. %s\n", i+1, step.Name)

			if step.Command != "" {
				fmt.Printf("   Command: %s\n", step.Command)
			}

			if step.Script != "" {
				fmt.Printf("   Script: %s\n", step.Script)
				if len(step.Args) > 0 {
					fmt.Printf("   Args: %v\n", step.Args)
				}
			}

			if step.WorkingDir != "" {
				fmt.Printf("   Working Directory: %s\n", step.WorkingDir)
			}

			if step.If != "" {
				fmt.Printf("   Condition: %s\n", step.If)
			}

			if step.Set != nil {
				fmt.Println("   Sets variables:")
				for k, v := range step.Set {
					fmt.Printf("     %s = %s\n", k, v)
				}
			}

			fmt.Println()
		}
	}

	return nil
}
