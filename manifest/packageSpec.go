package manifest

import manifest2 "github.com/thisismeamir/hepsw/internal/manifest"

// PackageSpec defines how to build and configure the package
type PackageSpec struct {
	Source       manifest2.SourceSpec        `yaml:"source"`
	Build        manifest2.BuildSpec         `yaml:"build"`
	Dependencies []manifest2.DependencySpec  `yaml:"dependencies,omitempty"`
	Options      []manifest2.OptionSpec      `yaml:"options,omitempty"`
	Environments []manifest2.EnvironmentSpec `yaml:"environments,omitempty"`
}
