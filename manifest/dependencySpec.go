package manifest

// DependencySpec defines a dependency with its requirements
type DependencySpec struct {
	Name        string            `yaml:"name"`
	Version     VersionConstraint `yaml:"version"`
	Optional    bool              `yaml:"optional,omitempty"`
	Conditional *ConditionalSpec  `yaml:"conditional,omitempty"`
	Options     map[string]string `yaml:"options,omitempty"`    // required options
	Targets     []string          `yaml:"targets,omitempty"`    // required targets
	BuildFlags  map[string]string `yaml:"buildFlags,omitempty"` // required flags
	Components  []string          `yaml:"components,omitempty"` // specific components needed
	Scope       DependencyScope   `yaml:"scope,omitempty"`
}
