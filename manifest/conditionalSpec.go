package manifest

// ConditionalSpec defines when a dependency is needed
type ConditionalSpec struct {
	When     string            `yaml:"when,omitempty"`     // option expression
	Platform []string          `yaml:"platform,omitempty"` // os constraints
	Arch     []string          `yaml:"arch,omitempty"`     // architecture constraints
	Compiler []string          `yaml:"compiler,omitempty"` // compiler constraints
	Options  map[string]string `yaml:"options,omitempty"`  // when these options are set
}
