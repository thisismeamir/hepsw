package manifestCmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
)

func runManifestNew(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Get template
	tmpl, err := getTemplate(manifestNewTemplate)
	if err != nil {
		return err
	}

	// Customize template with name
	tmpl.Name = name

	// Determine output path
	outputPath := manifestNewOutput
	if outputPath == "" {
		outputPath = fmt.Sprintf("%s.yaml", name)
	}

	// Save manifest
	if err := loader.SaveManifest(tmpl, outputPath); err != nil {
		return fmt.Errorf("failed to save manifest: %w", err)
	}

	fmt.Printf("Created manifest: %s\n", outputPath)
	fmt.Printf("Template: %s\n", manifestNewTemplate)
	return nil
}
