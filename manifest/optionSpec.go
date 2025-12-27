package manifest

// OptionSpec defines a build option/feature flag
type OptionSpec struct {
	Name        string        `yaml:"name"`
	Type        OptionType    `yaml:"type"`
	Description string        `yaml:"description,omitempty"`
	Default     interface{}   `yaml:"default,omitempty"`
	Allowed     []interface{} `yaml:"allowed,omitempty"`     // valid values
	Conflicts   []string      `yaml:"conflicts,omitempty"`   // conflicting options
	Requires    []string      `yaml:"requires,omitempty"`    // required options
	AffectsHash bool          `yaml:"affectsHash,omitempty"` // affects build hash
	Propagates  bool          `yaml:"propagates,omitempty"`  // propagates to dependents
}
