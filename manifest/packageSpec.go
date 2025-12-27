package manifest

// PackageSpec defines how to build and configure the package
type PackageSpec struct {
	Source       SourceSpec        `yaml:"source"`
	Build        BuildSpec         `yaml:"build"`
	Dependencies []DependencySpec  `yaml:"dependencies,omitempty"`
	Options      []OptionSpec      `yaml:"options,omitempty"`
	Environments []EnvironmentSpec `yaml:"environments,omitempty"`
}
