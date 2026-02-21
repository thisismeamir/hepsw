package manifestCmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
)

func runManifestCheck(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	fmt.Printf("Checking manifest: %s@%s\n", m.Name, m.Version)
	fmt.Println()

	// Run comprehensive validation
	result := manifest.ValidateManifest(m)

	// Additional deep checks
	deepIssues := performDeepChecks(m)

	// Print all issues
	hasIssues := false

	if len(result.Errors) > 0 {
		fmt.Println("ERRORS:")
		for _, err := range result.Errors {
			fmt.Printf("  ✗ %s\n", err)
		}
		fmt.Println()
		hasIssues = true
	}

	if len(result.Warnings) > 0 {
		fmt.Println("WARNINGS:")
		for _, warn := range result.Warnings {
			fmt.Printf("  ⚠ %s\n", warn)
		}
		fmt.Println()
		hasIssues = true
	}

	if len(deepIssues) > 0 {
		fmt.Println("CONSISTENCY ISSUES:")
		for _, issue := range deepIssues {
			fmt.Printf("  ⚠ %s\n", issue)
		}
		fmt.Println()
		hasIssues = true
	}

	if !hasIssues {
		fmt.Println("✓ All checks passed")
	}

	return nil
}

func performDeepChecks(m *manifest.Manifest) []string {
	issues := make([]string, 0)

	// Check for incompatible options
	// (This would require more complex logic to detect actual incompatibilities)

	// Check for circular dependencies in options
	// (Simplified check)

	// Check for missing toolchain for recipe steps
	hasCMake := false
	for _, tool := range m.Specifications.Build.Toolchain {
		if tool.Name == "cmake" {
			hasCMake = true
			break
		}
	}

	// Check if recipe uses cmake commands
	for _, step := range m.Recipe.Configuration {
		if strings.Contains(strings.ToLower(step.Command), "cmake") && !hasCMake {
			issues = append(issues, "Recipe uses cmake but cmake is not in toolchain")
			break
		}
	}

	return issues
}
