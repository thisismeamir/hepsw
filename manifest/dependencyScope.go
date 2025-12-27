package manifest

// DependencyScope defines when the dependency is needed
type DependencyScope string

const (
	ScopeBuild   DependencyScope = "build"   // needed at build time
	ScopeRuntime DependencyScope = "runtime" // needed at runtime
	ScopeTest    DependencyScope = "test"    // needed for testing
	ScopeAll     DependencyScope = "all"     // needed always (default)
)
