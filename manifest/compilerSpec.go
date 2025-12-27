package manifest

// CompilerSpec defines compiler requirements
type CompilerSpec struct {
	Type    string `yaml:"type"` // gcc, clang, msvc
	Version string `yaml:"version,omitempty"`
	Min     string `yaml:"min,omitempty"`
	Max     string `yaml:"max,omitempty"`
}
