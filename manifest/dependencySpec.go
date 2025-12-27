package manifest

import manifest2 "github.com/thisismeamir/hepsw/internal/manifest"

// DependencySpec defines a dependency with its requirements
type DependencySpec struct {
	Name        string                      `yaml:"name"`
	Version     manifest2.VersionConstraint `yaml:"version"`
	Optional    bool                        `yaml:"optional,omitempty"`
	Conditional *manifest2.ConditionalSpec  `yaml:"conditional,omitempty"`
	Options     map[string]string           `yaml:"options,omitempty"`    // required options
	Targets     []string                    `yaml:"targets,omitempty"`    // required targets
	BuildFlags  map[string]string           `yaml:"buildFlags,omitempty"` // required flags
	Components  []string                    `yaml:"components,omitempty"` // specific components needed
	Scope       manifest2.DependencyScope   `yaml:"scope,omitempty"`
}
