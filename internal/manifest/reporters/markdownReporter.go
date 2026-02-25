package reporters

import (
	"fmt"
	"strings"

	"github.com/thisismeamir/hepsw/internal/manifest"
)

func markdownHeaderAndBasicInfo(m *manifest.Manifest) string {
	var sb strings.Builder
	// Header
	sb.WriteString(fmt.Sprintf("# Manifest Report: %s\n\n", m.Name))

	// Basic Information
	sb.WriteString("## Basic Information\n\n")
	sb.WriteString(fmt.Sprintf("- **Name:** %s\n", m.Name))
	sb.WriteString(fmt.Sprintf("- **Version:** %s\n", m.Version))
	sb.WriteString(fmt.Sprintf("- **Description:** %s\n\n", m.Description))
	return sb.String()
}

func markdownSourceInfo(m *manifest.Manifest) string {
	var sb strings.Builder
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
	return sb.String()

}

func markdownMetadata(m *manifest.Manifest) string {
	var sb strings.Builder
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

	return sb.String()
}

func markdownBuildSpecifictations(m *manifest.Manifest) string {
	var sb strings.Builder
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

	return sb.String()
}

func markdownDependencies(m *manifest.Manifest) string {
	var sb strings.Builder
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
	return sb.String()
}

func generateMarkdownReport(m *manifest.Manifest) string {
	var sb strings.Builder
	headerAndBasicInfo := markdownHeaderAndBasicInfo(m)
	sb.WriteString(headerAndBasicInfo)

	sourceInfo := markdownSourceInfo(m)
	sb.WriteString(sourceInfo)

	metadata := markdownMetadata(m)
	sb.WriteString(metadata)

	buildSpecification := markdownBuildSpecifictations(m)
	sb.WriteString(buildSpecification)

	dependencies := markdownDependencies(m)
	sb.WriteString(dependencies)

	// Recipe
	sb.WriteString("## Recipe\n\n")
	printMarkdownRecipePhase(&sb, "Configuration", m.Recipe.Configuration)
	printMarkdownRecipePhase(&sb, "Build", m.Recipe.Build)
	printMarkdownRecipePhase(&sb, "Install", m.Recipe.Install)
	printMarkdownRecipePhase(&sb, "Use", m.Recipe.Use)

	return sb.String()
}

func printMarkdownRecipePhase(sb *strings.Builder, phaseName string, steps []manifest.RecipeStep) {
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
