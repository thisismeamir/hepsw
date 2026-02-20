package reporters

import (
	"fmt"
	"strings"

	"github.com/thisismeamir/hepsw/internal/manifest"
	"gopkg.in/yaml.v3"
)

// ReportFormat specifies the output format for reports
type ReportFormat string

const (
	FormatText     ReportFormat = "text"
	FormatMarkdown ReportFormat = "markdown"
	FormatJSON     ReportFormat = "json"
)

// GenerateReport creates a comprehensive report of a manifest
func GenerateReport(m *manifest.Manifest, format ReportFormat) (string, error) {
	switch format {
	case FormatText:
		return generateTextReport(m), nil
	case FormatMarkdown:
		return generateMarkdownReport(m), nil
	case FormatJSON:
		return generateJSONReport(m)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func generateJSONReport(m *manifest.Manifest) (string, error) {
	// Use existing YAML marshaling since manifest is already structured
	data, err := yaml.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("failed to marshal manifest: %w", err)
	}
	return string(data), nil
}

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
