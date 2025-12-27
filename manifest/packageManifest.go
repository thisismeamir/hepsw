package manifest

// PackageManifest represents the complete manifest file for a HEP software package
type PackageManifest struct {
	APIVersion string          `yaml:"apiVersion"` // e.g., "v1"
	Kind       string          `yaml:"kind"`       // Always "Package"
	Metadata   PackageMetadata `yaml:"metadata"`
	Spec       PackageSpec     `yaml:"spec"`
}
