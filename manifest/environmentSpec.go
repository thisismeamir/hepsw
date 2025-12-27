package manifest

// EnvironmentSpec defines environment-specific configurations
type EnvironmentSpec struct {
	Name      string            `yaml:"name"` // e.g., "linux-x86_64-gcc11"
	Platform  string            `yaml:"platform,omitempty"`
	Arch      string            `yaml:"arch,omitempty"`
	Compiler  CompilerSpec      `yaml:"compiler,omitempty"`
	Variables map[string]string `yaml:"variables,omitempty"`
	Overrides *SpecOverrides    `yaml:"overrides,omitempty"`
}
