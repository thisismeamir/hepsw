# HepSW Database Schema Documentation

## Overview

This database schema is designed for the HepSW package manager index. It stores metadata about packages, their versions, and dependencies, while the actual manifest files remain on GitHub.

## Schema Design Philosophy

1. **Metadata Registry**: The database is a discovery layer, not a storage system for manifests or source code
2. **Hybrid Dependencies**: Dependencies are stored as text names (always works) with optional FK references (when available)
3. **Immutable Versions**: Once published, version metadata doesn't change (except deprecation/yanked flags)
4. **Read-Heavy Optimization**: Indexed for common search and lookup patterns

## Tables

### 1. `packages`
Core package information (name, description, metadata).

**Key Constraints:**
- `name` must be unique
- Soft deletes not implemented (use `yanked` flag on versions instead)

### 2. `versions`
Each version of a package with its manifest location and source information.

**Key Constraints:**
- Unique constraint on `(package_id, version)`
- `source_type` must be 'git', 'tarball', or 'url'
- CASCADE delete when package is removed

**Important Fields:**
- `manifest_url`: Direct link to raw manifest file on GitHub
- `manifest_hash`: SHA256 for integrity verification
- `source_url`: Where to get the actual source code
- `source_ref`: Version-specific reference (git tag, commit, tarball checksum)

### 3. `dependencies`
Tracks what each version depends on.

**Hybrid Approach:**
- `dependency_name`: Always present (text string)
- `dependency_package_id`: Present only when dependency exists in our index
- This allows system dependencies (gcc, cmake) that aren't HepSW packages

## Common Queries

### Search for packages by name
```sql
-- Exact match
SELECT * FROM packages WHERE name = 'root';

-- Prefix search
SELECT * FROM packages WHERE name LIKE 'root%';

-- Fuzzy search (contains)
SELECT * FROM packages WHERE name LIKE '%root%';
```

### Get all versions of a package
```sql
SELECT v.version, v.manifest_url, v.published_at, v.deprecated, v.yanked
FROM versions v
JOIN packages p ON v.package_id = p.id
WHERE p.name = 'root'
ORDER BY v.published_at DESC;
```

### Get latest non-deprecated version
```sql
SELECT v.*
FROM versions v
JOIN packages p ON v.package_id = p.id
WHERE p.name = 'root'
  AND v.deprecated = 0
  AND v.yanked = 0
ORDER BY v.published_at DESC
LIMIT 1;

-- Or use the view:
SELECT * FROM latest_versions WHERE name = 'root';
```

### Get dependencies for a specific version
```sql
SELECT d.dependency_name, d.version_constraint, d.optional, d.condition
FROM dependencies d
JOIN versions v ON d.version_id = v.id
JOIN packages p ON v.package_id = p.id
WHERE p.name = 'root' AND v.version = '6.30.02';
```

### Reverse dependency lookup (what depends on X?)
```sql
-- What packages depend on cmake?
SELECT DISTINCT p.name, v.version, d.version_constraint
FROM dependencies d
JOIN versions v ON d.version_id = v.id
JOIN packages p ON v.package_id = p.id
WHERE d.dependency_name = 'cmake';

-- Or use the view:
SELECT * FROM reverse_dependencies WHERE dependency_name = 'cmake';
```

### Find packages by source repository
```sql
-- All packages from root-project
SELECT DISTINCT p.name, p.description
FROM packages p
JOIN versions v ON p.id = v.package_id
WHERE v.source_url LIKE '%github.com/root-project%';
```

### Search by tags
```sql
-- Find all HEP-related packages
SELECT * FROM packages 
WHERE tags LIKE '%hep%' OR description LIKE '%HEP%';
```

### Get package statistics
```sql
-- Use the pre-built view
SELECT * FROM package_stats ORDER BY version_count DESC;
```

## Data Integrity Rules

1. **Package Names**: Must be unique, lowercase recommended, no spaces
2. **Versions**: Follow semantic versioning when possible (not enforced)
3. **Manifest Hash**: Always SHA256, stored as hex string
4. **Source Types**: Only 'git', 'tarball', or 'url'
5. **Dependencies**: Name must match exactly for FK linking to work

## Migration Strategy

When adding this to Turso:

1. Create tables in order: `packages` → `versions` → `dependencies`
2. The triggers and views are created after tables
3. Initial population should:
   - First insert all packages
   - Then insert versions
   - Finally insert dependencies with FK linking where possible

## Example: Adding the ROOT package

```sql
-- 1. Insert package
INSERT INTO packages (name, description, documentation_url, maintainer, tags)
VALUES (
    'root',
    'CERN ROOT Data Analysis Framework',
    'https://root.cern/',
    'hepsw-team',
    'hep,analysis,cern,data-analysis'
);

-- 2. Insert version
INSERT INTO versions (
    package_id,
    version,
    manifest_url,
    manifest_hash,
    source_type,
    source_url,
    source_ref,
    notes
)
VALUES (
    (SELECT id FROM packages WHERE name = 'root'),
    '6.30.02',
    'https://raw.githubusercontent.com/HepSW/manifests/main/root/6.30.02.yaml',
    'abc123...', -- actual SHA256 hash
    'git',
    'https://github.com/root-project/root.git',
    'v6-30-02',
    'ROOT is the primary data analysis framework in HEP.'
);

-- 3. Insert dependencies
INSERT INTO dependencies (version_id, dependency_name, version_constraint, optional)
VALUES 
    ((SELECT id FROM versions WHERE version = '6.30.02' AND package_id = (SELECT id FROM packages WHERE name = 'root')),
     'cmake', '>=3.20', 0),
    ((SELECT id FROM versions WHERE version = '6.30.02' AND package_id = (SELECT id FROM packages WHERE name = 'root')),
     'gcc', '>=9.0', 0),
    ((SELECT id FROM versions WHERE version = '6.30.02' AND package_id = (SELECT id FROM packages WHERE name = 'root')),
     'python', '>=3.8', 0);
```

## Performance Considerations

- All foreign keys are indexed for JOIN performance
- `packages.name` has a unique index for fast lookups
- `versions(package_id, version)` has a composite unique index
- `dependencies.dependency_name` is indexed for reverse lookups
- Views are pre-defined for common complex queries

## Future Extensions

Possible additions without breaking changes:
- `package_aliases` table for alternative names
- `authors` table with many-to-many relationship
- `download_stats` table for popularity metrics
- `vulnerability_reports` table for security tracking
- Full-text search using SQLite FTS5 extension
