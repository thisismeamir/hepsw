package manifest

// BuildCommands defines custom build steps
type BuildCommands struct {
	PreConfigure []string `yaml:"preConfigure,omitempty"`
	Configure    []string `yaml:"configure,omitempty"`
	PreBuild     []string `yaml:"preBuild,omitempty"`
	Build        []string `yaml:"build,omitempty"`
	PostBuild    []string `yaml:"postBuild,omitempty"`
	Install      []string `yaml:"install,omitempty"`
	Test         []string `yaml:"test,omitempty"`
}
