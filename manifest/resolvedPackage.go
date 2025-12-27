package manifest

// ResolvedPackage represents a package with resolved options and dependencies
type ResolvedPackage struct {
	Manifest    *PackageManifest
	Options     map[string]interface{}
	Environment *EnvironmentSpec
	BuildHash   string // hash of package + options + dependencies
}
