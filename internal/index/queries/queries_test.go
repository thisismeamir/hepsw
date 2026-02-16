package queries

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create schema
	schema := `
	CREATE TABLE packages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		description TEXT NOT NULL,
		documentation_url TEXT,
		maintainer TEXT,
		tags TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE versions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		package_id INTEGER NOT NULL,
		version TEXT NOT NULL,
		manifest_url TEXT NOT NULL,
		manifest_hash TEXT NOT NULL,
		source_type TEXT NOT NULL,
		source_url TEXT NOT NULL,
		source_ref TEXT,
		notes TEXT,
		deprecated BOOLEAN NOT NULL DEFAULT 0,
		yanked BOOLEAN NOT NULL DEFAULT 0,
		published_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (package_id) REFERENCES packages(id) ON DELETE CASCADE,
		UNIQUE(package_id, version)
	);

	CREATE TABLE dependencies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		version_id INTEGER NOT NULL,
		dependency_name TEXT NOT NULL,
		dependency_package_id INTEGER,
		version_constraint TEXT NOT NULL,
		optional BOOLEAN NOT NULL DEFAULT 0,
		condition TEXT,
		FOREIGN KEY (version_id) REFERENCES versions(id) ON DELETE CASCADE,
		FOREIGN KEY (dependency_package_id) REFERENCES packages(id) ON DELETE SET NULL
	);
	`

	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	return db
}

// seedTestData inserts test data
func seedTestData(t *testing.T, db *sql.DB) {
	// Insert test package
	_, err := db.Exec(`
		INSERT INTO packages (name, description, documentation_url, maintainer, tags)
		VALUES ('root', 'CERN ROOT Data Analysis Framework', 'https://root.cern/', 'hepsw-team', 'hep,analysis,cern')
	`)
	if err != nil {
		t.Fatalf("Failed to insert test package: %v", err)
	}

	// Insert test version
	_, err = db.Exec(`
		INSERT INTO versions (package_id, version, manifest_url, manifest_hash, source_type, source_url, source_ref)
		VALUES (1, '6.30.02', 'https://example.com/root/6.30.02.yaml', 'abc123', 'git', 'https://github.com/root-project/root.git', 'v6-30-02')
	`)
	if err != nil {
		t.Fatalf("Failed to insert test version: %v", err)
	}

	// Insert test dependencies
	_, err = db.Exec(`
		INSERT INTO dependencies (version_id, dependency_name, version_constraint, optional)
		VALUES 
			(1, 'cmake', '>=3.20', 0),
			(1, 'gcc', '>=9.0', 0),
			(1, 'python', '>=3.8', 0)
	`)
	if err != nil {
		t.Fatalf("Failed to insert test dependencies: %v", err)
	}
}

func TestGetPackageByName(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedTestData(t, db)

	queries := New(db)
	ctx := context.Background()

	pkg, err := queries.GetPackageByName(ctx, "root")
	if err != nil {
		t.Fatalf("GetPackageByName failed: %v", err)
	}

	if pkg.Name != "root" {
		t.Errorf("Expected package name 'root', got '%s'", pkg.Name)
	}

	if pkg.Description != "CERN ROOT Data Analysis Framework" {
		t.Errorf("Unexpected description: %s", pkg.Description)
	}

	if pkg.Tags != "hep,analysis,cern" {
		t.Errorf("Unexpected tags: %s", pkg.Tags)
	}

	// Test non-existent package
	_, err = queries.GetPackageByName(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent package")
	}
}

func TestGetVersion(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedTestData(t, db)

	queries := New(db)
	ctx := context.Background()

	ver, err := queries.GetVersion(ctx, 1, "6.30.02")
	if err != nil {
		t.Fatalf("GetVersion failed: %v", err)
	}

	if ver.Version != "6.30.02" {
		t.Errorf("Expected version '6.30.02', got '%s'", ver.Version)
	}

	if ver.SourceType != "git" {
		t.Errorf("Expected source type 'git', got '%s'", ver.SourceType)
	}

	if ver.SourceURL != "https://github.com/root-project/root.git" {
		t.Errorf("Unexpected source URL: %s", ver.SourceURL)
	}
}

func TestGetDependencies(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedTestData(t, db)

	queries := New(db)
	ctx := context.Background()

	deps, err := queries.GetDependencies(ctx, 1)
	if err != nil {
		t.Fatalf("GetDependencies failed: %v", err)
	}

	if len(deps) != 3 {
		t.Errorf("Expected 3 dependencies, got %d", len(deps))
	}

	// Check first dependency
	if deps[0].DependencyName != "cmake" {
		t.Errorf("Expected first dependency 'cmake', got '%s'", deps[0].DependencyName)
	}

	if deps[0].VersionConstraint != ">=3.20" {
		t.Errorf("Expected version constraint '>=3.20', got '%s'", deps[0].VersionConstraint)
	}
}

func TestSearchPackages(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedTestData(t, db)

	queries := New(db)
	ctx := context.Background()

	// Exact match
	results, err := queries.SearchPackages(ctx, "root", true)
	if err != nil {
		t.Fatalf("SearchPackages failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	// Fuzzy search
	results, err = queries.SearchPackages(ctx, "root", false)
	if err != nil {
		t.Fatalf("SearchPackages failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
}

func TestGetLatestVersion(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedTestData(t, db)

	// Insert another version
	_, err := db.Exec(`
		INSERT INTO versions (package_id, version, manifest_url, manifest_hash, source_type, source_url, published_at)
		VALUES (1, '6.28.00', 'https://example.com/root/6.28.00.yaml', 'def456', 'git', 'https://github.com/root-project/root.git', ?)
	`, time.Now().Add(-24*time.Hour))
	if err != nil {
		t.Fatalf("Failed to insert second version: %v", err)
	}

	queries := New(db)
	ctx := context.Background()

	latest, err := queries.GetLatestVersion(ctx, 1)
	if err != nil {
		t.Fatalf("GetLatestVersion failed: %v", err)
	}

	// Should return 6.30.02 as it was inserted later
	if latest.Version != "6.30.02" {
		t.Errorf("Expected latest version '6.30.02', got '%s'", latest.Version)
	}
}

func TestSearchByTags(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedTestData(t, db)

	queries := New(db)
	ctx := context.Background()

	results, err := queries.SearchByTags(ctx, []string{"hep"})
	if err != nil {
		t.Fatalf("SearchByTags failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	if results[0].Name != "root" {
		t.Errorf("Expected package 'root', got '%s'", results[0].Name)
	}
}
