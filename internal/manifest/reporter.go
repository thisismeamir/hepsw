package manifest

import (
	"fmt"
	"strings"

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
func GenerateReport(m *Manifest, format ReportFormat) (string, error) {
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

func generateTextReport(m *Manifest) string {
	var sb strings.Builder

	// Header
	sb.WriteString(strings.Repeat("=", 80) + "\n")
	sb.WriteString(fmt.Sprintf("MANIFEST REPORT: %s\n", m.Name))
	sb.WriteString(strings.Repeat("=", 80) + "\n\n")

	// Basic Information
	sb.WriteString("BASIC INFORMATION\n")
	sb.WriteString(strings.Repeat("-", 80) + "\n")
	sb.WriteString(fmt.Sprintf("Name:        %s\n", m.Name))
	sb.WriteString(fmt.Sprintf("Version:     %s\n", m.Version))
	sb.WriteString(fmt.Sprintf("Description: %s\n\n", m.Description))

	// Source
	sb.WriteString("SOURCE\n")
	sb.WriteString(strings.Repeat("-", 80) + "\n")
	sb.WriteString(fmt.Sprintf("Type:     %s\n", m.Source.Type))
	sb.WriteString(fmt.Sprintf("URL:      %s\n", m.Source.Url))
	if m.Source.Tag != "" {
		sb.WriteString(fmt.Sprintf("Tag:      %s\n", m.Source.Tag))
	}
	if m.Source.Checksum != "" {
		sb.WriteString(fmt.Sprintf("Checksum: %s\n", m.Source.Checksum))
	}
	sb.WriteString("\n")

	// Metadata
	if len(m.Metadata.Authors) > 0 || m.Metadata.Homepage != "" || m.Metadata.License != "" {
		sb.WriteString("METADATA\n")
		sb.WriteString(strings.Repeat("-", 80) + "\n")
		if len(m.Metadata.Authors) > 0 {
			sb.WriteString("Authors:\n")
			for _, author := range m.Metadata.Authors {
				sb.WriteString(fmt.Sprintf("  - %s\n", author))
			}
		}
		if m.Metadata.Homepage != "" {
			sb.WriteString(fmt.Sprintf("Homepage:      %s\n", m.Metadata.Homepage))
		}
		if m.Metadata.License != "" {
			sb.WriteString(fmt.Sprintf("License:       %s\n", m.Metadata.License))
		}
		if m.Metadata.Documentation != "" {
			sb.WriteString(fmt.Sprintf("Documentation: %s\n", m.Metadata.Documentation))
		}
		sb.WriteString("\n")
	}

	// Build Specifications
	sb.WriteString("BUILD SPECIFICATIONS\n")
	sb.WriteString(strings.Repeat("-", 80) + "\n")

	if len(m.Specifications.Build.Toolchain) > 0 {
		sb.WriteString("Toolchain:\n")
		for _, tool := range m.Specifications.Build.Toolchain {
			sb.WriteString(fmt.Sprintf("  - %s %s\n", tool.Name, tool.Version))
		}
	}

	if len(m.Specifications.Build.Targets) > 0 {
		sb.WriteString("Targets:\n")
		for _, target := range m.Specifications.Build.Targets {
			sb.WriteString(fmt.Sprintf("  - %s (%s)\n", target.Name, target.Architecture))
		}
	}

	if len(m.Specifications.Build.Options) > 0 {
		sb.WriteString("Options:\n")
		for _, opt := range m.Specifications.Build.Options {
			sb.WriteString(fmt.Sprintf("  - %s\n", opt))
		}
	}
	sb.WriteString("\n")

	// Dependencies
	if len(m.Specifications.Build.Dependencies) > 0 || len(m.Specifications.Runtime.Dependencies) > 0 {
		sb.WriteString("DEPENDENCIES\n")
		sb.WriteString(strings.Repeat("-", 80) + "\n")

		if len(m.Specifications.Build.Dependencies) > 0 {
			sb.WriteString("Build Dependencies:\n")
			for _, dep := range m.Specifications.Build.Dependencies {
				sb.WriteString(fmt.Sprintf("  - %s %s", dep.Name, dep.Version))
				if len(dep.ForOptions) > 0 {
					sb.WriteString(fmt.Sprintf(" [for: %s]", strings.Join(dep.ForOptions, ", ")))
				}
				if dep.IsOptional {
					sb.WriteString(" (optional)")
				}
				sb.WriteString("\n")
			}
		}

		if len(m.Specifications.Runtime.Dependencies) > 0 {
			sb.WriteString("Runtime Dependencies:\n")
			for _, dep := range m.Specifications.Runtime.Dependencies {
				sb.WriteString(fmt.Sprintf("  - %s %s", dep.Name, dep.Version))
				if len(dep.ForOptions) > 0 {
					sb.WriteString(fmt.Sprintf(" [for: %s]", strings.Join(dep.ForOptions, ", ")))
				}
				sb.WriteString("\n")
			}
		}
		sb.WriteString("\n")
	}

	// Environment Variables
	if len(m.Specifications.Environment.Build) > 0 ||
		len(m.Specifications.Environment.Runtime) > 0 ||
		len(m.Specifications.Environment.Self) > 0 {
		sb.WriteString("ENVIRONMENT VARIABLES\n")
		sb.WriteString(strings.Repeat("-", 80) + "\n")

		if len(m.Specifications.Environment.Build) > 0 {
			sb.WriteString("Build Environment:\n")
			for k, v := range m.Specifications.Environment.Build {
				sb.WriteString(fmt.Sprintf("  %s = %s\n", k, v))
			}
		}

		if len(m.Specifications.Environment.Runtime) > 0 {
			sb.WriteString("Runtime Environment:\n")
			for k, v := range m.Specifications.Environment.Runtime {
				sb.WriteString(fmt.Sprintf("  %s = %s\n", k, v))
			}
		}

		if len(m.Specifications.Environment.Self) > 0 {
			sb.WriteString("Self Environment:\n")
			for k, v := range m.Specifications.Environment.Self {
				sb.WriteString(fmt.Sprintf("  %s = %s\n", k, v))
			}
		}
		sb.WriteString("\n")
	}

	// Recipe Summary
	sb.WriteString("RECIPE SUMMARY\n")
	sb.WriteString(strings.Repeat("-", 80) + "\n")
	sb.WriteString(fmt.Sprintf("Configuration steps: %d\n", len(m.Recipe.Configuration)))
	sb.WriteString(fmt.Sprintf("Build steps:         %d\n", len(m.Recipe.Build)))
	sb.WriteString(fmt.Sprintf("Install steps:       %d\n", len(m.Recipe.Install)))
	sb.WriteString(fmt.Sprintf("Use steps:           %d\n", len(m.Recipe.Use)))
	sb.WriteString("\n")

	// Recipe Details
	printRecipePhase(&sb, "Configuration", m.Recipe.Configuration)
	printRecipePhase(&sb, "Build", m.Recipe.Build)
	printRecipePhase(&sb, "Install", m.Recipe.Install)
	printRecipePhase(&sb, "Use", m.Recipe.Use)

	return sb.String()
}

func printRecipePhase(sb *strings.Builder, phaseName string, steps []RecipeStep) {
	if len(steps) == 0 {
		return
	}

	sb.WriteString(fmt.Sprintf("%s STEPS\n", strings.ToUpper(phaseName)))
	sb.WriteString(strings.Repeat("-", 80) + "\n")

	for i, step := range steps {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, step.Name))
		if step.Command != "" {
			sb.WriteString(fmt.Sprintf("   Command: %s\n", step.Command))
		}
		if step.Script != "" {
			sb.WriteString(fmt.Sprintf("   Script: %s\n", step.Script))
			if len(step.Args) > 0 {
				sb.WriteString(fmt.Sprintf("   Args: %v\n", step.Args))
			}
		}
		if step.WorkingDir != "" {
			sb.WriteString(fmt.Sprintf("   Working Dir: %s\n", step.WorkingDir))
		}
		if step.If != "" {
			sb.WriteString(fmt.Sprintf("   Condition: %s\n", step.If))
		}
		if step.Set != nil {
			sb.WriteString("   Sets variables:\n")
			for k, v := range step.Set {
				sb.WriteString(fmt.Sprintf("     %s = %s\n", k, v))
			}
		}
	}
	sb.WriteString("\n")
}

func generateMarkdownReport(m *Manifest) string {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("# Manifest Report: %s\n\n", m.Name))

	// Basic Information
	sb.WriteString("## Basic Information\n\n")
	sb.WriteString(fmt.Sprintf("- **Name:** %s\n", m.Name))
	sb.WriteString(fmt.Sprintf("- **Version:** %s\n", m.Version))
	sb.WriteString(fmt.Sprintf("- **Description:** %s\n\n", m.Description))

	// Source
	sb.WriteString("## Source\n\n")
	sb.WriteString(fmt.Sprintf("- **Type:** %s\n", m.Source.Type))
	sb.WriteString(fmt.Sprintf("- **URL:** %s\n", m.Source.Url))
	if m.Source.Tag != "" {
		sb.WriteString(fmt.Sprintf("- **Tag:** %s\n", m.Source.Tag))
	}
	if m.Source.Checksum != "" {
		sb.WriteString(fmt.Sprintf("- **Checksum:** `%s`\n", m.Source.Checksum))
	}
	sb.WriteString("\n")

	// Metadata
	if len(m.Metadata.Authors) > 0 || m.Metadata.Homepage != "" || m.Metadata.License != "" {
		sb.WriteString("## Metadata\n\n")
		if len(m.Metadata.Authors) > 0 {
			sb.WriteString("**Authors:**\n")
			for _, author := range m.Metadata.Authors {
				sb.WriteString(fmt.Sprintf("- %s\n", author))
			}
		}
		if m.Metadata.Homepage != "" {
			sb.WriteString(fmt.Sprintf("- **Homepage:** %s\n", m.Metadata.Homepage))
		}
		if m.Metadata.License != "" {
			sb.WriteString(fmt.Sprintf("- **License:** %s\n", m.Metadata.License))
		}
		if m.Metadata.Documentation != "" {
			sb.WriteString(fmt.Sprintf("- **Documentation:** %s\n", m.Metadata.Documentation))
		}
		sb.WriteString("\n")
	}

	// Build Specifications
	sb.WriteString("## Build Specifications\n\n")

	if len(m.Specifications.Build.Toolchain) > 0 {
		sb.WriteString("### Toolchain\n\n")
		for _, tool := range m.Specifications.Build.Toolchain {
			sb.WriteString(fmt.Sprintf("- %s %s\n", tool.Name, tool.Version))
		}
		sb.WriteString("\n")
	}

	if len(m.Specifications.Build.Targets) > 0 {
		sb.WriteString("### Targets\n\n")
		for _, target := range m.Specifications.Build.Targets {
			sb.WriteString(fmt.Sprintf("- %s (%s)\n", target.Name, target.Architecture))
		}
		sb.WriteString("\n")
	}

	if len(m.Specifications.Build.Options) > 0 {
		sb.WriteString("### Options\n\n")
		for _, opt := range m.Specifications.Build.Options {
			sb.WriteString(fmt.Sprintf("- `%s`\n", opt))
		}
		sb.WriteString("\n")
	}

	// Dependencies
	if len(m.Specifications.Build.Dependencies) > 0 || len(m.Specifications.Runtime.Dependencies) > 0 {
		sb.WriteString("## Dependencies\n\n")

		if len(m.Specifications.Build.Dependencies) > 0 {
			sb.WriteString("### Build Dependencies\n\n")
			sb.WriteString("| Package | Version | Options | Optional |\n")
			sb.WriteString("|---------|---------|---------|----------|\n")
			for _, dep := range m.Specifications.Build.Dependencies {
				opts := "-"
				if len(dep.ForOptions) > 0 {
					opts = strings.Join(dep.ForOptions, ", ")
				}
				optional := "-"
				if dep.IsOptional {
					optional = "Yes"
				}
				sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", dep.Name, dep.Version, opts, optional))
			}
			sb.WriteString("\n")
		}

		if len(m.Specifications.Runtime.Dependencies) > 0 {
			sb.WriteString("### Runtime Dependencies\n\n")
			sb.WriteString("| Package | Version | Options |\n")
			sb.WriteString("|---------|---------|----------|\n")
			for _, dep := range m.Specifications.Runtime.Dependencies {
				opts := "-"
				if len(dep.ForOptions) > 0 {
					opts = strings.Join(dep.ForOptions, ", ")
				}
				sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", dep.Name, dep.Version, opts))
			}
			sb.WriteString("\n")
		}
	}

	// Recipe
	sb.WriteString("## Recipe\n\n")
	printMarkdownRecipePhase(&sb, "Configuration", m.Recipe.Configuration)
	printMarkdownRecipePhase(&sb, "Build", m.Recipe.Build)
	printMarkdownRecipePhase(&sb, "Install", m.Recipe.Install)
	printMarkdownRecipePhase(&sb, "Use", m.Recipe.Use)

	return sb.String()
}

func printMarkdownRecipePhase(sb *strings.Builder, phaseName string, steps []RecipeStep) {
	if len(steps) == 0 {
		return
	}

	sb.WriteString(fmt.Sprintf("### %s Steps\n\n", phaseName))

	for i, step := range steps {
		sb.WriteString(fmt.Sprintf("%d. **%s**\n", i+1, step.Name))
		if step.Command != "" {
			sb.WriteString(fmt.Sprintf("   ```bash\n   %s\n   ```\n", step.Command))
		}
		if step.Script != "" {
			sb.WriteString(fmt.Sprintf("   - Script: `%s`\n", step.Script))
			if len(step.Args) > 0 {
				sb.WriteString(fmt.Sprintf("   - Args: `%v`\n", step.Args))
			}
		}
		if step.If != "" {
			sb.WriteString(fmt.Sprintf("   - Condition: `%s`\n", step.If))
		}
		sb.WriteString("\n")
	}
}

func generateJSONReport(m *Manifest) (string, error) {
	// Use existing YAML marshaling since manifest is already structured
	data, err := yaml.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("failed to marshal manifest: %w", err)
	}
	return string(data), nil
}

// GenerateDependencyReport creates a dependency-focused report
func GenerateDependencyReport(m *Manifest) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Dependency Report: %s@%s\n", m.Name, m.Version))
	sb.WriteString(strings.Repeat("=", 80) + "\n\n")

	allDeps := make(map[string]Dependency)

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
