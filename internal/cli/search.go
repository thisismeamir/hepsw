package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thisismeamir/hepsw/internal/index"
)

func Search(ctx context.Context, idx *index.Index, name string) {
	packages, err := idx.SearchPackages(ctx, name)
	if err != nil {
		PrintError("search failed: " + err.Error())
	}

	if len(packages) == 0 {
		PrintInfo("No packages found")
		return
	}

	PrintSection("Found " + string(len(packages)) + " packages:\n")
	for i, pkg := range packages {
		fmt.Printf("%i. %s:\n", i, pkg.Name)
		fmt.Printf("   %s\n", pkg.Description)
		if pkg.Tags != "" {
			PrintInfo("   Tags: %s\n" + pkg.Tags)
		}
		fmt.Println()
	}
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
	fmt.Printf("Created: %s\n", pkg.CreatedAt.Format("2006-01-02"))

	if ver != nil {
		PrintBullet("\nLatest Version: %s\n", ver.Version)
		PrintBullet("Published: %s\n", ver.PublishedAt.Format("2006-01-02"))
		PrintBullet("Manifest: %s\n", ver.ManifestURL)
		PrintBullet("Source: %s (%s)\n", ver.SourceURL, ver.SourceType)
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

	PrintSection("Versions of ", name, ":\n\n")
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

		PrintSection("Dependencies for %s@%s:\n\n", name, version)
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
		PrintInfo("No packages depend on ", name, "\n")
		return
	}

	PrintSection("Packages that depend on ", name, ":\n\n")
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
