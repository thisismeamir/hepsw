package queries

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/thisismeamir/hepsw/internal/index/models"
)

// Queries handles all database queries
type Queries struct {
	db *sql.DB
}

// New creates a new Queries instance
func New(db *sql.DB) *Queries {
	return &Queries{db: db}
}

// GetPackageByName retrieves a package by its name
func (q *Queries) GetPackageByName(ctx context.Context, name string) (*models.Package, error) {
	query := `
		SELECT id, name, description, documentation_url, maintainer, tags, created_time, updated_time
		FROM packages
		WHERE name = ?
	`

	var pkg models.Package
	err := q.db.QueryRowContext(ctx, query, name).Scan(
		&pkg.ID,
		&pkg.Name,
		&pkg.Description,
		&pkg.DocumentationURL,
		&pkg.Maintainer,
		&pkg.Tags,
		&pkg.CreatedTime,
		&pkg.UpdatedTime,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("package '%s' not found", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query package: %w", err)
	}

	return &pkg, nil
}

// GetPackageByID retrieves a package by its ID
func (q *Queries) GetPackageByID(ctx context.Context, id int64) (*models.Package, error) {
	query := `
		SELECT id, name, description, documentation_url, maintainer, tags, created_time, updated_time
		FROM packages
		WHERE id = ?
	`

	var pkg models.Package
	err := q.db.QueryRowContext(ctx, query, id).Scan(
		&pkg.ID,
		&pkg.Name,
		&pkg.Description,
		&pkg.DocumentationURL,
		&pkg.Maintainer,
		&pkg.Tags,
		&pkg.CreatedTime,
		&pkg.UpdatedTime,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("package with id %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query package: %w", err)
	}

	return &pkg, nil
}

// SearchPackages searches for packages by name (prefix or contains)
func (q *Queries) SearchPackages(ctx context.Context, searchTerm string, exactMatch bool) ([]models.Package, error) {
	var query string
	var args []interface{}

	if exactMatch {
		query = `
			SELECT id, name, description, documentation_url, maintainer, tags, created_time, updated_time
			FROM packages
			WHERE name = ?
			ORDER BY name
		`
		args = []interface{}{searchTerm}
	} else {
		query = `
			SELECT id, name, description, documentation_url, maintainer, tags, created_time, updated_time
			FROM packages
			WHERE name LIKE ?
			ORDER BY name
		`
		args = []interface{}{"%" + searchTerm + "%"}
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search packages: %w", err)
	}
	defer rows.Close()

	var packages []models.Package
	for rows.Next() {
		var pkg models.Package
		err := rows.Scan(
			&pkg.ID,
			&pkg.Name,
			&pkg.Description,
			&pkg.DocumentationURL,
			&pkg.Maintainer,
			&pkg.Tags,
			&pkg.CreatedTime,
			&pkg.UpdatedTime,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan package: %w", err)
		}
		packages = append(packages, pkg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating packages: %w", err)
	}

	return packages, nil
}

// ListPackages retrieves all packages with optional limit and offset
func (q *Queries) ListPackages(ctx context.Context, limit, offset int) ([]models.Package, error) {
	query := `
		SELECT id, name, description, documentation_url, maintainer, tags, created_time, updated_time
		FROM packages
		ORDER BY name
		LIMIT ? OFFSET ?
	`

	rows, err := q.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list packages: %w", err)
	}
	defer rows.Close()

	var packages []models.Package
	for rows.Next() {
		var pkg models.Package
		err := rows.Scan(
			&pkg.ID,
			&pkg.Name,
			&pkg.Description,
			&pkg.DocumentationURL,
			&pkg.Maintainer,
			&pkg.Tags,
			&pkg.CreatedTime,
			&pkg.UpdatedTime,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan package: %w", err)
		}
		packages = append(packages, pkg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating packages: %w", err)
	}

	return packages, nil
}

// GetVersionsByPackage retrieves all versions for a package
func (q *Queries) GetVersionsByPackage(ctx context.Context, packageID int64) ([]models.Version, error) {
	query := `
		SELECT id, package_id, version, manifest_url, manifest_hash, source_type, 
		       source_url, source_ref, notes, deprecated, yanked, published_at
		FROM versions
		WHERE package_id = ?
		ORDER BY published_at DESC
	`

	rows, err := q.db.QueryContext(ctx, query, packageID)
	if err != nil {
		return nil, fmt.Errorf("failed to query versions: %w", err)
	}
	defer rows.Close()

	var versions []models.Version
	for rows.Next() {
		var v models.Version
		err := rows.Scan(
			&v.ID,
			&v.PackageID,
			&v.Version,
			&v.ManifestURL,
			&v.ManifestHash,
			&v.SourceType,
			&v.SourceURL,
			&v.SourceRef,
			&v.Notes,
			&v.Deprecated,
			&v.Yanked,
			&v.PublishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan version: %w", err)
		}
		versions = append(versions, v)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating versions: %w", err)
	}

	return versions, nil
}

// GetVersion retrieves a specific version of a package
func (q *Queries) GetVersion(ctx context.Context, packageID int64, version string) (*models.Version, error) {
	query := `
		SELECT id, package_id, version, manifest_url, manifest_hash, source_type,
		       source_url, source_ref, notes, deprecated, yanked, published_at
		FROM versions
		WHERE package_id = ? AND version = ?
	`

	var v models.Version
	err := q.db.QueryRowContext(ctx, query, packageID, version).Scan(
		&v.ID,
		&v.PackageID,
		&v.Version,
		&v.ManifestURL,
		&v.ManifestHash,
		&v.SourceType,
		&v.SourceURL,
		&v.SourceRef,
		&v.Notes,
		&v.Deprecated,
		&v.Yanked,
		&v.PublishedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("version '%s' not found", version)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query version: %w", err)
	}

	return &v, nil
}

// GetLatestVersion retrieves the latest non-deprecated, non-yanked version
func (q *Queries) GetLatestVersion(ctx context.Context, packageID int64) (*models.Version, error) {
	query := `
		SELECT id, package_id, version, manifest_url, manifest_hash, source_type,
		       source_url, source_ref, notes, deprecated, yanked, published_at
		FROM versions
		WHERE package_id = ? AND deprecated = 0 AND yanked = 0
		ORDER BY published_at DESC
		LIMIT 1
	`

	var v models.Version
	err := q.db.QueryRowContext(ctx, query, packageID).Scan(
		&v.ID,
		&v.PackageID,
		&v.Version,
		&v.ManifestURL,
		&v.ManifestHash,
		&v.SourceType,
		&v.SourceURL,
		&v.SourceRef,
		&v.Notes,
		&v.Deprecated,
		&v.Yanked,
		&v.PublishedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no available version found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query latest version: %w", err)
	}

	return &v, nil
}

// GetDependencies retrieves all dependencies for a version
func (q *Queries) GetDependencies(ctx context.Context, versionID int64) ([]models.Dependency, error) {
	query := `
		SELECT id, version_id, dependency_name, dependency_package_id, 
		       version_constraint, optional, condition
		FROM dependencies
		WHERE version_id = ?
		ORDER BY optional, dependency_name
	`

	rows, err := q.db.QueryContext(ctx, query, versionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query dependencies: %w", err)
	}
	defer rows.Close()

	var deps []models.Dependency
	for rows.Next() {
		var d models.Dependency
		err := rows.Scan(
			&d.ID,
			&d.VersionID,
			&d.DependencyName,
			&d.DependencyPackageID,
			&d.VersionConstraint,
			&d.Optional,
			&d.Condition,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan dependency: %w", err)
		}
		deps = append(deps, d)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating dependencies: %w", err)
	}

	return deps, nil
}

// GetReverseDependencies finds what depends on a given package
func (q *Queries) GetReverseDependencies(ctx context.Context, packageName string) ([]models.ReverseDependency, error) {
	query := `
		SELECT d.dependency_name, p.name, v.version, d.version_constraint, d.optional
		FROM dependencies d
		JOIN versions v ON d.version_id = v.id
		JOIN packages p ON v.package_id = p.id
		WHERE d.dependency_name = ?
		ORDER BY p.name, v.version
	`

	rows, err := q.db.QueryContext(ctx, query, packageName)
	if err != nil {
		return nil, fmt.Errorf("failed to query reverse dependencies: %w", err)
	}
	defer rows.Close()

	var revDeps []models.ReverseDependency
	for rows.Next() {
		var rd models.ReverseDependency
		err := rows.Scan(
			&rd.DependencyName,
			&rd.DependentPackage,
			&rd.DependentVersion,
			&rd.VersionConstraint,
			&rd.Optional,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reverse dependency: %w", err)
		}
		revDeps = append(revDeps, rd)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reverse dependencies: %w", err)
	}

	return revDeps, nil
}

// GetPackageStats retrieves statistics for all packages
func (q *Queries) GetPackageStats(ctx context.Context) ([]models.PackageStats, error) {
	query := `
		SELECT p.id, p.name, p.description, COUNT(v.id) as version_count, MAX(v.published_at) as latest_release
		FROM packages p
		LEFT JOIN versions v ON p.id = v.package_id
		GROUP BY p.id, p.name, p.description
		ORDER BY p.name
	`

	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query package stats: %w", err)
	}
	defer rows.Close()

	var stats []models.PackageStats
	for rows.Next() {
		var s models.PackageStats
		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Description,
			&s.VersionCount,
			&s.LatestRelease,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan package stats: %w", err)
		}
		stats = append(stats, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating package stats: %w", err)
	}

	return stats, nil
}

// SearchByTags finds packages matching any of the given tags
func (q *Queries) SearchByTags(ctx context.Context, tags []string) ([]models.Package, error) {
	if len(tags) == 0 {
		return []models.Package{}, nil
	}

	// Build query with OR conditions for each tag
	query := `
		SELECT DISTINCT id, name, description, documentation_url, maintainer, tags, created_time, updated_time
		FROM packages
		WHERE `

	args := make([]interface{}, 0, len(tags))
	for i, tag := range tags {
		if i > 0 {
			query += " OR "
		}
		query += "tags LIKE ?"
		args = append(args, "%"+tag+"%")
	}
	query += " ORDER BY name"

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search by tags: %w", err)
	}
	defer rows.Close()

	var packages []models.Package
	for rows.Next() {
		var pkg models.Package
		err := rows.Scan(
			&pkg.ID,
			&pkg.Name,
			&pkg.Description,
			&pkg.DocumentationURL,
			&pkg.Maintainer,
			&pkg.Tags,
			&pkg.CreatedTime,
			&pkg.UpdatedTime,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan package: %w", err)
		}
		packages = append(packages, pkg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating packages: %w", err)
	}

	return packages, nil
}
