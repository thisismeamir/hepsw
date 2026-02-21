package manifestCmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
	"github.com/thisismeamir/hepsw/internal/manifest/reporters"
)

func runManifestShow(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	// Generate report based on format
	var output string
	switch showFormat {
	case "text":
		output, err = reporters.GenerateReport(m, reporters.FormatText)
	case "yaml":
		output, err = reporters.GenerateReport(m, reporters.FormatJSON) // Will output YAML
	case "json":
		output, err = reporters.GenerateReport(m, reporters.FormatJSON)
	case "markdown", "md":
		output, err = reporters.GenerateReport(m, reporters.FormatMarkdown)
	default:
		return fmt.Errorf("unsupported format: %s", showFormat)
	}

	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	fmt.Print(output)
	return nil
}
