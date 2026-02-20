package reporters

import (
	"fmt"
	"strings"

	"github.com/thisismeamir/hepsw/internal/manifest"
)

// GenerateDependencyReport creates a dependency-focused report
func GenerateDependencyReport(m *manifest.Manifest) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Dependency Report: %s@%s\n", m.Name, m.Version))
	sb.WriteString(strings.Repeat("=", 80) + "\n\n")

	allDeps := make(map[string]manifest.Dependency)

	// Collect all dependencies
	for _, dep := range m.Specifications.Build.Dependencies {
		key := fmt.Sprintf("%s (build)", dep.Name)
		allDeps[key] = dep
	}

	for _, dep := range m.Specifications.Runtime.Dependencies {
		key := fmt.Sprintf("%s (runtime)", dep.Name)
		allDeps[key] = dep
	}

	if len(allDeps) == 0 {
		sb.WriteString("No dependencies found.\n")
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("Total Dependencies: %d\n\n", len(allDeps)))

	for key, dep := range allDeps {
		sb.WriteString(fmt.Sprintf("- %s\n", key))
		sb.WriteString(fmt.Sprintf("  Version: %s\n", dep.Version))
		if len(dep.ForOptions) > 0 {
			sb.WriteString(fmt.Sprintf("  Required for options: %s\n", strings.Join(dep.ForOptions, ", ")))
		}
		if len(dep.WithOptions) > 0 {
			sb.WriteString(fmt.Sprintf("  Requires options: %s\n", strings.Join(dep.WithOptions, ", ")))
		}
		if dep.IsOptional {
			sb.WriteString("  Optional: Yes\n")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
