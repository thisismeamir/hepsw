package resolver

import (
	"context"
	"fmt"

	"github.com/thisismeamir/hepsw/internal/index/models"
	"github.com/thisismeamir/hepsw/internal/index/queries"
)

// DependencyNode represents a node in the dependency tree
type DependencyNode struct {
	Package      string            `json:"package"`
	Version      string            `json:"version"`
	Constraint   string            `json:"constraint,omitempty"`
	Optional     bool              `json:"optional"`
	Dependencies []*DependencyNode `json:"dependencies,omitempty"`
	Depth        int               `json:"depth"`
}

// Resolver handles dependency resolution
type Resolver struct {
	queries *queries.Queries
}

// New creates a new dependency resolver
func New(q *queries.Queries) *Resolver {
	return &Resolver{queries: q}
}

// ResolveDependencies resolves the full dependency tree for a package version
func (r *Resolver) ResolveDependencies(ctx context.Context, packageName, version string, includeOptional bool) (*DependencyNode, error) {
	// Get the package
	pkg, err := r.queries.GetPackageByName(ctx, packageName)
	if err != nil {
		return nil, fmt.Errorf("failed to get package: %w", err)
	}

	// Get the specific version
	var ver *models.Version
	if version == "latest" || version == "" {
		ver, err = r.queries.GetLatestVersion(ctx, pkg.ID)
	} else {
		ver, err = r.queries.GetVersion(ctx, pkg.ID, version)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
	}

	// Build the dependency tree
	root := &DependencyNode{
		Package: packageName,
		Version: ver.Version,
		Depth:   0,
	}

	visited := make(map[string]bool)
	if err := r.buildDependencyTree(ctx, root, ver.ID, visited, includeOptional, 0); err != nil {
		return nil, err
	}

	return root, nil
}

// buildDependencyTree recursively builds the dependency tree
func (r *Resolver) buildDependencyTree(
	ctx context.Context,
	node *DependencyNode,
	versionID int64,
	visited map[string]bool,
	includeOptional bool,
	depth int,
) error {
	// Prevent infinite recursion
	if depth > 100 {
		return fmt.Errorf("dependency depth limit exceeded (possible circular dependency)")
	}

	// Get dependencies for this version
	deps, err := r.queries.GetDependencies(ctx, versionID)
	if err != nil {
		return fmt.Errorf("failed to get dependencies: %w", err)
	}

	for _, dep := range deps {
		// Skip optional dependencies if not requested
		if dep.Optional && !includeOptional {
			continue
		}

		depKey := fmt.Sprintf("%s@%s", dep.DependencyName, dep.VersionConstraint)

		// Check if we've already visited this dependency
		if visited[depKey] {
			// Add reference but don't recurse
			node.Dependencies = append(node.Dependencies, &DependencyNode{
				Package:    dep.DependencyName,
				Constraint: dep.VersionConstraint,
				Optional:   dep.Optional,
				Depth:      depth + 1,
			})
			continue
		}

		visited[depKey] = true

		// Try to resolve the dependency if it exists in our index
		depNode := &DependencyNode{
			Package:    dep.DependencyName,
			Constraint: dep.VersionConstraint,
			Optional:   dep.Optional,
			Depth:      depth + 1,
		}

		// If the dependency has a package ID, we can recurse
		if dep.DependencyPackageID != nil {
			depPkg, err := r.queries.GetPackageByID(ctx, *dep.DependencyPackageID)
			if err == nil {
				// Get latest version of dependency
				depVer, err := r.queries.GetLatestVersion(ctx, depPkg.ID)
				if err == nil {
					depNode.Version = depVer.Version
					// Recursively resolve dependencies
					if err := r.buildDependencyTree(ctx, depNode, depVer.ID, visited, includeOptional, depth+1); err != nil {
						return err
					}
				}
			}
		}

		node.Dependencies = append(node.Dependencies, depNode)
	}

	return nil
}

// FlattenDependencies returns a flat list of all dependencies
func (r *Resolver) FlattenDependencies(root *DependencyNode) []string {
	seen := make(map[string]bool)
	var result []string

	var flatten func(*DependencyNode)
	flatten = func(node *DependencyNode) {
		for _, dep := range node.Dependencies {
			key := fmt.Sprintf("%s@%s", dep.Package, dep.Constraint)
			if !seen[key] {
				seen[key] = true
				result = append(result, dep.Package)
				flatten(dep)
			}
		}
	}

	flatten(root)
	return result
}

// GetAllDependencies returns both direct and transitive dependencies
func (r *Resolver) GetAllDependencies(ctx context.Context, packageName, version string) ([]string, error) {
	tree, err := r.ResolveDependencies(ctx, packageName, version, false)
	if err != nil {
		return nil, err
	}
	return r.FlattenDependencies(tree), nil
}

// CheckCircularDependencies checks if there are circular dependencies
func (r *Resolver) CheckCircularDependencies(ctx context.Context, packageName, version string) (bool, []string, error) {
	pkg, err := r.queries.GetPackageByName(ctx, packageName)
	if err != nil {
		return false, nil, err
	}

	var ver *models.Version
	if version == "latest" || version == "" {
		ver, err = r.queries.GetLatestVersion(ctx, pkg.ID)
	} else {
		ver, err = r.queries.GetVersion(ctx, pkg.ID, version)
	}
	if err != nil {
		return false, nil, err
	}

	path := []string{packageName}
	visited := make(map[string]bool)
	inPath := make(map[string]bool)
	inPath[packageName] = true

	hasCircular, circularPath := r.detectCycle(ctx, ver.ID, visited, inPath, path)
	return hasCircular, circularPath, nil
}

// detectCycle performs DFS to detect circular dependencies
func (r *Resolver) detectCycle(
	ctx context.Context,
	versionID int64,
	visited map[string]bool,
	inPath map[string]bool,
	path []string,
) (bool, []string) {
	deps, err := r.queries.GetDependencies(ctx, versionID)
	if err != nil {
		return false, nil
	}

	for _, dep := range deps {
		if dep.Optional {
			continue
		}

		depName := dep.DependencyName

		// If we've seen this in our current path, we have a cycle
		if inPath[depName] {
			return true, append(path, depName)
		}

		// If already visited in a previous path, skip
		if visited[depName] {
			continue
		}

		// If the dependency exists in our index, recurse
		if dep.DependencyPackageID != nil {
			depPkg, err := r.queries.GetPackageByID(ctx, *dep.DependencyPackageID)
			if err != nil {
				continue
			}

			depVer, err := r.queries.GetLatestVersion(ctx, depPkg.ID)
			if err != nil {
				continue
			}

			visited[depName] = true
			inPath[depName] = true
			newPath := append(path, depName)

			if hasCycle, cyclePath := r.detectCycle(ctx, depVer.ID, visited, inPath, newPath); hasCycle {
				return true, cyclePath
			}

			delete(inPath, depName)
		}
	}

	return false, nil
}
