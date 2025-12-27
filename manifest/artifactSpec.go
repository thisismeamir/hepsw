package manifest

// ArtifactSpec defines what to keep after building
type ArtifactSpec struct {
	Type    string   `yaml:"type"` // binary, library, header, data
	Paths   []string `yaml:"paths"`
	Include []string `yaml:"include,omitempty"` // glob patterns
	Exclude []string `yaml:"exclude,omitempty"`
}
