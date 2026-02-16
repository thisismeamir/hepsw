package index

import (
	"context"
	"fmt"

	"github.com/thisismeamir/hepsw/internal/index/client"
	"github.com/thisismeamir/hepsw/internal/index/models"
	"github.com/thisismeamir/hepsw/internal/index/resolver"
)

// Index is the main interface for interacting with the HepSW package index
type Index struct {
	client   *client.Client
	resolver *resolver.Resolver
}

// New creates a new Index instance
func New(config *client.IndexConfig) (*Index, error) {
	c, err := client.New(config)
	if err != nil {
		return nil, err
	}

	return &Index{
		client:   c,
		resolver: resolver.New(c.Queries()),
	}, nil
}

// Close closes the index connection
func (idx *Index) Close() error {
	return idx.client.Close()
}

// ============================================================================
// Package Operations
// ============================================================================

// GetPackage retrieves a package by name
func (idx *Index) GetPackage(ctx context.Context, name string) (*models.Package, error) {
	// Check cache first
	if idx.client.Cache() != nil {
		if cached, found := idx.client.Cache().Get("pkg:" + name); found {
			return cached.(*models.Package), nil
		}
	}

	pkg, err := idx.client.Queries().GetPackageByName(ctx, name)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if idx.client.Cache() != nil {
		idx.client.Cache().Set("pkg:"+name, pkg)
	}

	return pkg, nil
}

// SearchPackages searches for packages by name
func (idx *Index) SearchPackages(ctx context.Context, query string) ([]models.Package, error) {
	return idx.client.Queries().SearchPackages(ctx, query, false)
}

// ListPackages lists all packages with pagination
func (idx *Index) ListPackages(ctx context.Context, limit, offset int) ([]models.Package, error) {
	return idx.client.Queries().ListPackages(ctx, limit, offset)
}

// SearchByTags finds packages by tags
func (idx *Index) SearchByTags(ctx context.Context, tags []string) ([]models.Package, error) {
	return idx.client.Queries().SearchByTags(ctx, tags)
}

// ============================================================================
// Version Operations
// ============================================================================

// GetVersion retrieves a specific version of a package
func (idx *Index) GetVersion(ctx context.Context, packageName, version string) (*models.Version, error) {
	cacheKey := fmt.Sprintf("ver:%s:%s", packageName, version)

	// Check cache
	if idx.client.Cache() != nil {
		if cached, found := idx.client.Cache().Get(cacheKey); found {
			return cached.(*models.Version), nil
		}
	}

	pkg, err := idx.GetPackage(ctx, packageName)
	if err != nil {
		return nil, err
	}

	var ver *models.Version
	if version == "latest" || version == "" {
		ver, err = idx.client.Queries().GetLatestVersion(ctx, pkg.ID)
	} else {
		ver, err = idx.client.Queries().GetVersion(ctx, pkg.ID, version)
	}

	if err != nil {
		return nil, err
	}

	// Cache the result
	if idx.client.Cache() != nil {
		idx.client.Cache().Set(cacheKey, ver)
	}

	return ver, nil
}

// GetLatestVersion retrieves the latest version of a package
func (idx *Index) GetLatestVersion(ctx context.Context, packageName string) (*models.Version, error) {
	return idx.GetVersion(ctx, packageName, "latest")
}

// GetAllVersions retrieves all versions of a package
func (idx *Index) GetAllVersions(ctx context.Context, packageName string) ([]models.Version, error) {
	pkg, err := idx.GetPackage(ctx, packageName)
	if err != nil {
		return nil, err
	}

	return idx.client.Queries().GetVersionsByPackage(ctx, pkg.ID)
}

// ============================================================================
// Dependency Operations
// ============================================================================

// GetDependencies retrieves direct dependencies for a package version
func (idx *Index) GetDependencies(ctx context.Context, packageName, version string) ([]models.Dependency, error) {
	ver, err := idx.GetVersion(ctx, packageName, version)
	if err != nil {
		return nil, err
	}

	return idx.client.Queries().GetDependencies(ctx, ver.ID)
}

// GetReverseDependencies finds what depends on a package
func (idx *Index) GetReverseDependencies(ctx context.Context, packageName string) ([]models.ReverseDependency, error) {
	return idx.client.Queries().GetReverseDependencies(ctx, packageName)
}

// ResolveDependencyTree resolves the full dependency tree
func (idx *Index) ResolveDependencyTree(ctx context.Context, packageName, version string, includeOptional bool) (*resolver.DependencyNode, error) {
	return idx.resolver.ResolveDependencies(ctx, packageName, version, includeOptional)
}

// GetAllDependencies returns a flat list of all dependencies (direct and transitive)
func (idx *Index) GetAllDependencies(ctx context.Context, packageName, version string) ([]string, error) {
	return idx.resolver.GetAllDependencies(ctx, packageName, version)
}

// CheckCircularDependencies checks for circular dependencies
func (idx *Index) CheckCircularDependencies(ctx context.Context, packageName, version string) (bool, []string, error) {
	return idx.resolver.CheckCircularDependencies(ctx, packageName, version)
}

// ============================================================================
// Statistics and Discovery
// ============================================================================

// GetPackageStats retrieves statistics for all packages
func (idx *Index) GetPackageStats(ctx context.Context) ([]models.PackageStats, error) {
	return idx.client.Queries().GetPackageStats(ctx)
}

// ============================================================================
// Utility Methods
// ============================================================================

// Ping checks if the connection is alive
func (idx *Index) Ping(ctx context.Context) error {
	return idx.client.Ping(ctx)
}

// ClearCache clears the local cache
func (idx *Index) ClearCache() {
	if idx.client.Cache() != nil {
		idx.client.Cache().Clear()
	}
}
