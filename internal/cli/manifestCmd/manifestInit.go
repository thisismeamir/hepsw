package manifestCmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
)

func runManifestInit(cmd *cobra.Command, args []string) error {
	fmt.Println("Initializing new HepSW manifest...")
	fmt.Println()

	// Interactive prompts
	name := prompt("Package name", getDefaultNameOrNameOfDir())
	version := prompt("Version", "0.1.0")
	description := prompt("Description", "")

	sourceType := promptChoice("Source type", []string{"git", "tarball", "svn", "local"}, "git")
	sourceURL := prompt("Source URL", "")

	var sourceTag string
	if sourceType == "git" {
		sourceTag = prompt("Git tag/branch", "main")
	}

	buildSystem := promptChoice("Build system", []string{"cmake", "autotools", "make", "custom"}, "cmake")

	// Create manifest
	m := &manifest.Manifest{
		Name:        name,
		Version:     version,
		Description: description,
		Source: manifest.SourceSpec{
			Type: sourceType,
			Url:  sourceURL,
			Tag:  sourceTag,
		},
		Specifications: manifest.Specifications{
			Build: manifest.BuildSpecification{
				Toolchain: []manifest.Tool{},
			},
		},
		Recipe: createRecipeForBuildSystem(buildSystem),
	}

	// Save manifest
	outputPath := fmt.Sprintf("%s.yaml", name)
	if err := loader.SaveManifest(m, outputPath); err != nil {
		return fmt.Errorf("failed to save manifest: %w", err)
	}

	fmt.Printf("\nCreated manifest: %s\n", outputPath)
	return nil
}
