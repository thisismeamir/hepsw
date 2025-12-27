package manifest

// VersionConstraint defines acceptable version range
type VersionConstraint struct {
	Exact      string   `yaml:"exact,omitempty"`      // "1.2.3"
	Min        string   `yaml:"min,omitempty"`        // "1.0.0"
	Max        string   `yaml:"max,omitempty"`        // "2.0.0"
	Range      string   `yaml:"range,omitempty"`      // ">=1.0.0,<2.0.0"
	Compatible string   `yaml:"compatible,omitempty"` // "^1.2.3" or "~1.2.3"
	Exclude    []string `yaml:"exclude,omitempty"`    // versions to exclude
}
