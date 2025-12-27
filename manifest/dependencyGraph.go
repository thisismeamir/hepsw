package manifest

// DependencyGraph represents the resolved dependency tree
type DependencyGraph struct {
	Root         *ResolvedPackage
	Dependencies map[string]*ResolvedPackage // keyed by name+version+hash
	BuildOrder   []string                    // topologically sorted
}
