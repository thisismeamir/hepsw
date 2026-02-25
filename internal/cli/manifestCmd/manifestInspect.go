package manifestCmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
	"github.com/thisismeamir/hepsw/internal/utils"
)

func runManifestInspect(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]
	field := ""
	if len(args) > 1 {
		field = args[1]
	}

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	accessor := manifest.NewManifestAccessor(m)

	// If no field specified, show summary
	if field == "" {
		fmt.Printf("Manifest: %s@%s\n", accessor.Name(), accessor.Version())
		fmt.Printf("Description: %s\n\n", accessor.Description())

		fmt.Println("Available fields:")
		fmt.Println("  name, version, description")
		fmt.Println("  source, source.type, source.url, source.tag, source.checksum")
		fmt.Println("  metadata, metadata.authors, metadata.license, metadata.homepage")
		fmt.Println("  specifications.build.dependencies")
		fmt.Println("  specifications.runtime.dependencies")
		fmt.Println("  specifications.build.toolchain")
		fmt.Println("  specifications.build.targets")
		fmt.Println("  specifications.build.options")
		fmt.Println("  recipe, recipe.configuration, recipe.build, recipe.install")
		return nil
	}

	// Inspect specific field
	value, err := inspectField(accessor, field)
	if err != nil {
		return err
	}

	fmt.Println(value)
	return nil
}

func inspectField(accessor *manifest.ManifestAccessor, field string) (string, error) {
	switch field {
	case "name":
		return accessor.Name(), nil
	case "version":
		return accessor.Version(), nil
	case "description":
		return accessor.Description(), nil
	case "source":
		return fmt.Sprintf("Type: %s\nURL: %s\nTag: %s",
			accessor.SourceType(), accessor.SourceURL(), accessor.SourceTag()), nil
	case "source.type":
		return accessor.SourceType(), nil
	case "source.url":
		return accessor.SourceURL(), nil
	case "source.tag":
		return accessor.SourceTag(), nil
	case "source.checksum":
		return accessor.SourceChecksum(), nil
	case "metadata.authors":
		return utils.FormatList(accessor.Authors()), nil
	case "metadata.license":
		return accessor.License(), nil
	case "metadata.homepage":
		return accessor.Homepage(), nil
	case "specifications.build.dependencies":
		return utils.FormatDependencies(accessor.BuildDependencies()), nil
	case "specifications.runtime.dependencies":
		return utils.FormatDependencies(accessor.RuntimeDependencies()), nil
	case "specifications.build.toolchain":
		return utils.FormatToolChain(accessor.Toolchain()), nil
	case "specifications.build.options":
		return utils.FormatList(accessor.Options()), nil
	default:
		return "", fmt.Errorf("unknown field: %s", field)
	}
}
