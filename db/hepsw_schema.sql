-- HepSW Package Index Database Schema
-- SQLite/Turso compatible

-- ============================================================================
-- Packages Table
-- ============================================================================
CREATE TABLE packages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    documentation_url TEXT,
    maintainer TEXT,
    tags TEXT, -- Comma-separated tags: "hep,analysis,cern"
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for name search (exact and prefix matching)
CREATE INDEX idx_packages_name ON packages(name);

-- Index for full-text search on description
CREATE INDEX idx_packages_description ON packages(description);

-- ============================================================================
-- Versions Table
-- ============================================================================
CREATE TABLE versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    package_id INTEGER NOT NULL,
    version TEXT NOT NULL,
    manifest_url TEXT NOT NULL,
    manifest_hash TEXT NOT NULL, -- SHA256 hash
    source_type TEXT NOT NULL CHECK(source_type IN ('git', 'tarball', 'url')),
    source_url TEXT NOT NULL,
    source_ref TEXT, -- tag/commit/checksum depending on source_type
    notes TEXT,
    deprecated BOOLEAN NOT NULL DEFAULT 0,
    yanked BOOLEAN NOT NULL DEFAULT 0,
    published_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (package_id) REFERENCES packages(id) ON DELETE CASCADE,
    UNIQUE(package_id, version) -- A package can't have duplicate versions
);

-- Index for querying versions by package
CREATE INDEX idx_versions_package_id ON versions(package_id);

-- Index for source URL discovery queries
CREATE INDEX idx_versions_source_url ON versions(source_url);

-- Index for finding non-yanked, non-deprecated versions quickly
CREATE INDEX idx_versions_status ON versions(deprecated, yanked);

-- ============================================================================
-- Dependencies Table (Hybrid approach)
-- ============================================================================
CREATE TABLE dependencies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    version_id INTEGER NOT NULL,
    dependency_name TEXT NOT NULL, -- Always store the name
    dependency_package_id INTEGER, -- Nullable FK when it exists in our index
    version_constraint TEXT NOT NULL, -- e.g., ">=3.20", "^1.0.0", "*"
    optional BOOLEAN NOT NULL DEFAULT 0,
    condition TEXT, -- Description of when this dependency is needed
    
    FOREIGN KEY (version_id) REFERENCES versions(id) ON DELETE CASCADE,
    FOREIGN KEY (dependency_package_id) REFERENCES packages(id) ON DELETE SET NULL
);

-- Index for forward dependency lookup (what does this version depend on?)
CREATE INDEX idx_dependencies_version_id ON dependencies(version_id);

-- Index for reverse dependency lookup (what depends on this package?)
CREATE INDEX idx_dependencies_dependency_name ON dependencies(dependency_name);

-- Index for reverse lookup via FK when available
CREATE INDEX idx_dependencies_dependency_package_id ON dependencies(dependency_package_id);

-- ============================================================================
-- Triggers for automatic timestamp updates
-- ============================================================================
CREATE TRIGGER update_packages_timestamp 
AFTER UPDATE ON packages
FOR EACH ROW
BEGIN
    UPDATE packages SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

-- ============================================================================
-- Views for common queries
-- ============================================================================

-- View: Latest version of each package (non-yanked, non-deprecated)
CREATE VIEW latest_versions AS
SELECT 
    p.name,
    p.description,
    v.version,
    v.manifest_url,
    v.published_at
FROM packages p
JOIN versions v ON p.id = v.package_id
WHERE v.deprecated = 0 AND v.yanked = 0
AND v.published_at = (
    SELECT MAX(v2.published_at)
    FROM versions v2
    WHERE v2.package_id = v.package_id
    AND v2.deprecated = 0
    AND v2.yanked = 0
);

-- View: Package with all version count
CREATE VIEW package_stats AS
SELECT 
    p.id,
    p.name,
    p.description,
    COUNT(v.id) as version_count,
    MAX(v.published_at) as latest_release
FROM packages p
LEFT JOIN versions v ON p.id = v.package_id
GROUP BY p.id, p.name, p.description;

-- View: Reverse dependencies (what depends on what)
CREATE VIEW reverse_dependencies AS
SELECT 
    d.dependency_name,
    p.name as dependent_package,
    v.version as dependent_version,
    d.version_constraint,
    d.optional
FROM dependencies d
JOIN versions v ON d.version_id = v.id
JOIN packages p ON v.package_id = p.id
ORDER BY d.dependency_name, p.name;
