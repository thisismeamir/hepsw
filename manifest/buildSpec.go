package manifest

import manifest2 "github.com/thisismeamir/hepsw/internal/manifest"

// BuildSpec defines how to build the package
type BuildSpec struct {
	System    string                   `yaml:"system"`              // cmake, autotools, make, custom
	Directory string                   `yaml:"directory,omitempty"` // subdirectory to build in
	Commands  manifest2.BuildCommands  `yaml:"commands,omitempty"`
	Flags     map[string]string        `yaml:"flags,omitempty"`    // CMAKE_CXX_STANDARD: "17"
	Targets   []string                 `yaml:"targets,omitempty"`  // specific build targets
	Parallel  *int                     `yaml:"parallel,omitempty"` // number of parallel jobs
	Timeout   *int                     `yaml:"timeout,omitempty"`  // build timeout in seconds
	Artifacts []manifest2.ArtifactSpec `yaml:"artifacts,omitempty"`
}
