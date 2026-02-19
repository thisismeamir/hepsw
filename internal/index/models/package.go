package models

import "time"

// Package represents a software package in the HepSW index
type Package struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	DocumentationURL *string   `json:"documentation_url,omitempty"`
	Maintainer       *string   `json:"maintainer,omitempty"`
	Tags             string    `json:"tags"` // Comma-separated
	CreatedTime      time.Time `json:"created_time"`
	UpdatedTime      time.Time `json:"updated_time"`
}

// GetTags returns the tags as a slice
func (p *Package) GetTags() []string {
	if p.Tags == "" {
		return []string{}
	}
	tags := []string{}
	current := ""
	for _, c := range p.Tags {
		if c == ',' {
			if current != "" {
				tags = append(tags, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		tags = append(tags, current)
	}
	return tags
}

// Version represents a specific version of a package
type Version struct {
	ID           int64     `json:"id"`
	PackageID    int64     `json:"package_id"`
	Version      string    `json:"version"`
	ManifestURL  string    `json:"manifest_url"`
	ManifestHash string    `json:"manifest_hash"` // SHA256
	SourceType   string    `json:"source_type"`   // git, tarball, url
	SourceURL    string    `json:"source_url"`
	SourceRef    *string   `json:"source_ref,omitempty"`
	Notes        *string   `json:"notes,omitempty"`
	Deprecated   bool      `json:"deprecated"`
	Yanked       bool      `json:"yanked"`
	PublishedAt  time.Time `json:"published_at"`
}

// IsAvailable returns true if the version is not deprecated or yanked
func (v *Version) IsAvailable() bool {
	return !v.Deprecated && !v.Yanked
}

// Dependency represents a dependency relationship
type Dependency struct {
	ID                  int64   `json:"id"`
	VersionID           int64   `json:"version_id"`
	DependencyName      string  `json:"dependency_name"`
	DependencyPackageID *int64  `json:"dependency_package_id,omitempty"`
	VersionConstraint   string  `json:"version_constraint"`
	Optional            bool    `json:"optional"`
	Condition           *string `json:"condition,omitempty"`
}

// PackageWithVersions combines package info with all its versions
type PackageWithVersions struct {
	Package  Package   `json:"package"`
	Versions []Version `json:"versions"`
}

// VersionWithDependencies combines version info with its dependencies
type VersionWithDependencies struct {
	Version      Version      `json:"version"`
	Dependencies []Dependency `json:"dependencies"`
}

// PackageStats represents statistics about a package
type PackageStats struct {
	ID            int64      `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	VersionCount  int        `json:"version_count"`
	LatestRelease *time.Time `json:"latest_release,omitempty"`
}

// LatestVersion represents the latest version view from database
type LatestVersion struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	ManifestURL string    `json:"manifest_url"`
	PublishedAt time.Time `json:"published_at"`
}

// ReverseDependency represents what depends on what
type ReverseDependency struct {
	DependencyName    string `json:"dependency_name"`
	DependentPackage  string `json:"dependent_package"`
	DependentVersion  string `json:"dependent_version"`
	VersionConstraint string `json:"version_constraint"`
	Optional          bool   `json:"optional"`
}
