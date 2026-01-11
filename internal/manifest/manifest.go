package manifest

type Manifest struct {
	Name           string           `yaml:"name"`
	Version        string           `yaml:"version"`
	Description    string           `yaml:"description"`
	Source         SourceSpec       `yaml:"source"`
	Metadata       ManifestMetaData `yaml:"metadata,omitempty"`
	Specifications Specifications   `yaml:"specifications"`
	Recipe         Recipe           `yaml:"recipe"`
}

type SourceSpec struct {
	Type     string `yaml:"type"`
	Url      string `yaml:"url"`
	Tag      string `yaml:"tag,omitempty"`
	Checksum string `yaml:"checksum,omitempty"`
}

type ManifestMetaData struct {
	Authors       []string `yaml:"authors"`
	Homepage      string   `yaml:"homepage"`
	License       string   `yaml:"license"`
	Documentation string   `yaml:"documentation"`
}

type Specifications struct {
	Build       BuildSpecification       `yaml:"build"`
	Runtime     RuntimeSpecification     `yaml:"runtime"`
	Environment EnvironmentSpecification `yaml:"environment"`
}

type EnvironmentSpecification struct {
	Build   map[string]string `yaml:"build"`
	Runtime map[string]string `yaml:"runtime"`
	Self    map[string]string `yaml:"self"`
}

type RuntimeSpecification struct {
	Dependencies []Dependency `yaml:"dependencies"`
}

type BuildSpecification struct {
	Toolchain    []Tool              `yaml:"toolchain"`
	Targets      []Targets           `yaml:"targets"`
	Options      []string            `yaml:"options"`
	Dependencies []Dependency        `yaml:"dependencies"`
	Variables    []map[string]string `yaml:"variables"`
}

type Dependency struct {
	Name        string   `yaml:"name"`
	Version     string   `yaml:"version"`
	ForOptions  []string `yaml:"forOptions"`
	WithOptions []string `yaml:"withOptions"`
	isOptional  bool     `yaml:"isOptional"`
}

type Tool struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type Targets struct {
	Name         string `yaml:"name"`
	Architecture string `yaml:"architecture"`
}

type Recipe struct {
	Configuration []RecipeStep `yaml:"configuration"`
	Build         []RecipeStep `yaml:"build"`
	Install       []RecipeStep `yaml:"install"`
	Use           []RecipeStep `yaml:"use"`
}

type RecipeStep struct {
	Name       string            `yaml:"name"`
	Command    string            `yaml:"command,omitempty"`
	Script     string            `yaml:"script,omitempty"`
	Args       []string          `yaml:"args,omitempty"`
	WorkingDir string            `yaml:"working_dir,omitempty"`
	If         string            `yaml:"if,omitempty"`
	Set        map[string]string `yaml:"set,omitempty"`
}
