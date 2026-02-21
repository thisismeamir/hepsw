package manifestCmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
	"github.com/thisismeamir/hepsw/internal/manifest/reporters"
	"gopkg.in/yaml.v3"
)

func runManifestExport(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	var output string

	switch exportFormat {
	case "json":
		data, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		output = string(data)

	case "yaml":
		data, err := yaml.Marshal(m)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML: %w", err)
		}
		output = string(data)

	case "lockfile":
		output = generateLockfile(m)

	case "report":
		output, err = reporters.GenerateReport(m, reporters.FormatText)
		if err != nil {
			return fmt.Errorf("failed to generate report: %w", err)
		}

	default:
		return fmt.Errorf("unsupported export format: %s", exportFormat)
	}

	// Write output
	if exportOutput != "" {
		if err := os.WriteFile(exportOutput, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		fmt.Printf("Exported to: %s\n", exportOutput)
	} else {
		fmt.Print(output)
	}

	return nil
}

func generateLockfile(m *manifest.Manifest) string {
	var sb strings.Builder

	sb.WriteString("# HepSW Lockfile\n")
	sb.WriteString(fmt.Sprintf("# Generated for: %s@%s\n\n", m.Name, m.Version))

	sb.WriteString("[[package]]\n")
	sb.WriteString(fmt.Sprintf("name = \"%s\"\n", m.Name))
	sb.WriteString(fmt.Sprintf("version = \"%s\"\n", m.Version))
	sb.WriteString(fmt.Sprintf("source = \"%s\"\n", m.Source.Url))
	if m.Source.Checksum != "" {
		sb.WriteString(fmt.Sprintf("checksum = \"%s\"\n", m.Source.Checksum))
	}
	sb.WriteString("\n")

	// Add dependencies
	for _, dep := range m.Specifications.Build.Dependencies {
		sb.WriteString("[[dependency]]\n")
		sb.WriteString(fmt.Sprintf("name = \"%s\"\n", dep.Name))
		sb.WriteString(fmt.Sprintf("version = \"%s\"\n", dep.Version))
		sb.WriteString("type = \"build\"\n\n")
	}

	for _, dep := range m.Specifications.Runtime.Dependencies {
		sb.WriteString("[[dependency]]\n")
		sb.WriteString(fmt.Sprintf("name = \"%s\"\n", dep.Name))
		sb.WriteString(fmt.Sprintf("version = \"%s\"\n", dep.Version))
		sb.WriteString("type = \"runtime\"\n\n")
	}

	return sb.String()
}
