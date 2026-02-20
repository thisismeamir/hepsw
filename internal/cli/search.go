package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/configuration"
	"github.com/thisismeamir/hepsw/internal/index"
	"github.com/thisismeamir/hepsw/internal/index/models"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Searches a package agains the local database.",
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
				Search(context.Background(), newIndex, searchPackageIdentity.Name)
			}
		}
	},
}

func Search(ctx context.Context, idx *index.Index, name string) {
	packages, err := idx.SearchPackages(ctx, name)
	if err != nil {
		PrintError("search failed: " + err.Error())
	}

	if len(packages) == 0 {
		PrintInfo("No packages found")
	} else {
		PrintSection("Found packages for " + name + ":")
		for i, pkg := range packages {
			fmt.Printf("%d. %s:\n", i, colorInfo(pkg.Name))
			fmt.Printf("   %s\n", pkg.Description)
		}
	}
}

func Versions(ctx context.Context, idx *index.Index, name string) {
	versions, err := idx.GetAllVersions(ctx, name)
	if err != nil {
		PrintError("Failed to get versions: " + err.Error())
	}

	if len(versions) == 0 {
		PrintInfo("No versions found")
		return
	}

	PrintSection("Versions of " + name + ":\n\n")
	for _, v := range versions {
		status := ""
		if v.Deprecated {
			status = " [DEPRECATED]"
		}
		if v.Yanked {
			status = " [YANKED]"
		}
		PrintBullet("  %s%s (published: %s)\n", v.Version, status, v.PublishedAt.Format("2006-01-02"))
	}
}

func Dependencies(ctx context.Context, idx *index.Index, name, version string, showTree bool) {
	if showTree {
		tree, err := idx.ResolveDependencyTree(ctx, name, version, false)
		if err != nil {
			PrintError("Failed to resolve dependencies: " + err.Error())
		}

		data, err := json.MarshalIndent(tree, "", "  ")
		if err != nil {
			PrintError("Failed to marshal tree: " + err.Error())
		}
		fmt.Println(string(data))
	} else {
		deps, err := idx.GetDependencies(ctx, name, version)
		if err != nil {
			PrintError("Failed to get dependencies: " + err.Error())
		}

		if len(deps) == 0 {
			PrintInfo("No dependencies")
			return
		}

		PrintSection("Dependencies for %s@%s:\n\n" + name + version)
		for _, dep := range deps {
			optional := ""
			if dep.Optional {
				optional = " (optional)"
			}
			PrintBullet(fmt.Sprintf("%s %s%s\n", dep.DependencyName, dep.VersionConstraint, optional))
			if dep.Condition != nil {
				fmt.Printf("    Condition: %s\n", *dep.Condition)
			}
		}
	}
}

func ReverseDependency(ctx context.Context, idx *index.Index, name string) {
	revDeps, err := idx.GetReverseDependencies(ctx, name)
	if err != nil {
		PrintError("Failed to get reverse dependencies: %v" + err.Error())
	}

	if len(revDeps) == 0 {
		PrintInfo("No packages depend on " + name + "\n")
		return
	}

	PrintSection("Packages that depend on " + name + ":\n\n")
	for _, rd := range revDeps {
		optional := ""
		if rd.Optional {
			optional = " (optional)"
		}
		point := fmt.Sprintf("%s@%s requires %s%s\n",
			rd.DependentPackage,
			rd.DependentVersion,
			rd.VersionConstraint,
			optional)
		PrintBullet(point)
	}
}

func List(ctx context.Context, idx *index.Index) {
	packages, err := idx.ListPackages(ctx, 100, 0)
	if err != nil {
		PrintError("Failed to list packages: " + err.Error())
	}

	totalPackages := fmt.Sprintf("Total packages: %d\n\n", len(packages))
	PrintSection(totalPackages)
	for _, pkg := range packages {
		desc := pkg.Description
		if len(desc) > 60 {
			desc = desc[:57] + "..."
		}
		list := fmt.Sprintf("%-20s %s\n", pkg.Name, desc)
		PrintBullet(list)
	}
}

func Stats(ctx context.Context, idx *index.Index) {
	stats, err := idx.GetPackageStats(ctx)
	if err != nil {
		PrintError("Failed to get stats: " + err.Error())
	}

	PrintSection("Package Statistics:\n")

	totalVersions := 0
	for _, s := range stats {
		totalVersions += s.VersionCount
	}

	totalPackages := fmt.Sprintf("Total packages: %d\n", len(stats))
	totalVErsions := fmt.Sprintf("Total versions: %d\n", totalVersions)

	PrintInfo(totalPackages)
	PrintInfo(totalVErsions)

	// Sort by version count
	PrintSection("Top packages by version count:")
	for i, s := range stats {
		if i >= 10 {
			break
		}
		releaseInfo := "no releases"
		if s.LatestRelease != nil {
			releaseInfo = s.LatestRelease.Format("2006-01-02")
		}

		// Truncate description
		desc := s.Description
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}

		fmt.Printf("  %2d. %-20s (%2d versions, latest: %s)\n",
			i+1, s.Name, s.VersionCount, releaseInfo)
		fmt.Printf("      %s\n", desc)
	}
}
