package index

type DependencyNode struct {
	Package      string
	Version      string
	Constraint   string
	Optional     bool
	Dependencies []*DependencyNode
	Depth        int
}

type Resolver struct {
	queries *Query
}
