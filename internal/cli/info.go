package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/configuration"
	"github.com/thisismeamir/hepsw/internal/index"
	"github.com/thisismeamir/hepsw/internal/index/models"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Shows detailed information about a package.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			PrintError("Package name is required.")
		}
		if len(args) >= 1 {
			for _, pkgName := range args {
				searchPackageIdentity, err := models.GetSearchPackageIdentity(pkgName)
				if err != nil {
					PrintError(err.Error())
					panic(err)
				}
				config, err := configuration.GetConfiguration()
				if err != nil {
					PrintError(err.Error())
					panic(err)
				}
				newIndex, indexError := index.New(&config.IndexConfig)
				if indexError != nil {
					PrintError(indexError.Error())
					panic(indexError)
				}
				Info(context.Background(), newIndex, searchPackageIdentity.Name)

			}
		}
	},
}

func Info(ctx context.Context, idx *index.Index, name string) {
	pkg, err := idx.GetPackage(ctx, name)
	if err != nil {
		PrintError("Failed to get package: " + err.Error())
	}

	ver, err := idx.GetLatestVersion(ctx, name)
	if err != nil {
		PrintWarning("Warning: Could not get latest version: " + err.Error())
	}

	PrintSection("Package: %s\n" + pkg.Name)
	fmt.Printf("Description: %s\n", pkg.Description)
	if pkg.DocumentationURL != nil {
		PrintBullet("Documentation: %s\n", *pkg.DocumentationURL)
	}
	if pkg.Maintainer != nil {
		PrintBullet("Maintainer: %s\n", *pkg.Maintainer)
	}
	if pkg.Tags != "" {
		PrintBullet("Tags: %s\n", pkg.Tags)
	}
	fmt.Printf("Created: %s\n", pkg.CreatedTime.Format("2006-01-02"))

	if ver != nil {
		PrintBullet("\nLatest Version: %s\n", ver.Version)
		PrintBullet("Published: %s\n", ver.PublishedAt.Format("2006-01-02"))
		PrintBullet("Manifest: %s\n", ver.ManifestURL)
		PrintBullet("Source: %s (%s)\n", ver.SourceURL, ver.SourceType)
	}
}
