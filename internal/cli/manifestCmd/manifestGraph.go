package manifestCmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
)

func runManifestGraph(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	// Generate graph
	graph := generateDependencyGraph(m)

	fmt.Println(graph)

	return nil
}

func generateDependencyGraph(m *manifest.Manifest) string {
	var sb strings.Builder

	sb.WriteString("# Dependency Graph\n\n")
	sb.WriteString(fmt.Sprintf("Package: %s@%s\n", m.Name, m.Version))
	sb.WriteString("\n")

	// Build dependencies
	sb.WriteString("Build Dependencies:\n")
	if len(m.Specifications.Build.Dependencies) > 0 {
		for _, dep := range m.Specifications.Build.Dependencies {
			sb.WriteString(fmt.Sprintf("  %s --> %s %s", m.Name, dep.Name, dep.Version))
			if len(dep.ForOptions) > 0 {
				sb.WriteString(fmt.Sprintf(" [when: %s]", strings.Join(dep.ForOptions, ", ")))
			}
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString("  (none)\n")
	}

	sb.WriteString("\n")

	// Runtime dependencies
	sb.WriteString("Runtime Dependencies:\n")
	if len(m.Specifications.Runtime.Dependencies) > 0 {
		for _, dep := range m.Specifications.Runtime.Dependencies {
			sb.WriteString(fmt.Sprintf("  %s ==> %s %s", m.Name, dep.Name, dep.Version))
			if len(dep.ForOptions) > 0 {
				sb.WriteString(fmt.Sprintf(" [when: %s]", strings.Join(dep.ForOptions, ", ")))
			}
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString("  (none)\n")
	}

	sb.WriteString("\n")

	// Recipe flow
	sb.WriteString("Recipe Flow:\n")
	sb.WriteString("  Configuration --> Build --> Install --> Use\n")

	return sb.String()
}
