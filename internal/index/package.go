package index

import (
	"time"
)

type Package struct {
	ID               int64
	Name             string
	Description      string
	DocumentationURL *string
	Maintainer       *string
	Tags             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Version struct {
	ID           int64
	PackageID    int64
	Version      string
	ManifestURL  string
	ManifestHash string
	SourceType   string
	SourceURL    string
	SourceRef    *string
	Notes        *string
	Deprecated   bool
	Yanked       bool
	PublishedAt  time.Time
}

type Dependency struct {
	ID                  int64
	VersionID           int64
	DependencyName      string
	DependencyPackageID *int64
	VersionConstraint   string
	Optional            bool
	Condition           *string
}

type PackageWithVersions struct {
	Package  *Package
	Versions []*Version
}

type VersionWithDependencies struct {
	Version      *Version
	Dependencies []*Dependency
}

type PackageInfo struct {
	ID            int64
	Name          string
	Description   string
	VersionCount  int
	LatestRelease string
}

type LatestVersion struct {
	Name        string
	Description string
	Version     *Version
	ManifestURL string
	PublishedAt time.Time
}

type ReverseDependency struct {
	DependencyName     string
	DependentPackageID *int64
	DependencyVersion  *int64
	VersionConstraint  string
	Optional           bool
}
