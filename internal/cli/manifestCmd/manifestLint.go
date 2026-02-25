package manifestCmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
)

func runManifestLint(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	fmt.Printf("Linting manifest: %s@%s\n", m.Name, m.Version)
	fmt.Println()

	// Run validation (includes many lint checks)
	result := manifest.ValidateManifest(m)

	// Additional lint-specific checks
	lintIssues := performLintChecks(m)

	// Combine results
	allIssues := append(result.Warnings, lintIssues...)

	if len(allIssues) > 0 {
		fmt.Println("Lint Issues:")
		for _, issue := range allIssues {
			fmt.Printf("  ⚠ %s\n", issue)
		}
		fmt.Printf("\nFound %d issue(s)\n", len(allIssues))
	} else {
		fmt.Println("✓ No lint issues found")
	}

	return nil
}

func performLintChecks(m *manifest.Manifest) []manifest.ValidationError {
	issues := make([]manifest.ValidationError, 0)

	// Check for overly long description
	if len(m.Description) > 200 {
		issues = append(issues, manifest.ValidationError{
			Field:    "description",
			Message:  "Description is very long (>200 chars), consider shortening",
			Severity: "warning",
		})
	}

	// Check for missing homepage in metadata
	if m.Metadata.Homepage == "" {
		issues = append(issues, manifest.ValidationError{
			Field:    "metadata.homepage",
			Message:  "Homepage URL is recommended for discoverability",
			Severity: "info",
		})
	}

	// Check for empty recipe phases
	if len(m.Recipe.Build) == 0 {
		issues = append(issues, manifest.ValidationError{
			Field:    "recipe.build",
			Message:  "Build phase has no steps",
			Severity: "warning",
		})
	}

	return issues
}
