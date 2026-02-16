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
