package manifest

// SpecOverrides allows environment-specific overrides
type SpecOverrides struct {
	Build        *BuildSpec             `yaml:"build,omitempty"`
	Dependencies []DependencySpec       `yaml:"dependencies,omitempty"`
	Options      map[string]interface{} `yaml:"options,omitempty"`
}
