package configuration

type Configuration struct {
	Workspace     string `yaml:"workspace"`
	SourcesDir    string `yaml:"sourcesDir"`
	BuildsDir     string `yaml:"buildsDir"`
	InstallsDir   string `yaml:"installsDir"`
	EnvsDir       string `yaml:"envsDir"`
	LogsDir       string `yaml:"logsDir"`
	ToolchainsDir string `yaml:"toolchainsDir"`
	indexDir      string `yaml:"indexDir"`

	State      WorkspaceState `yaml:"state"`
	UserConfig UserConfig     `yaml:"userConfig"`
}

type WorkspaceState struct {
	Packages     map[string]interface{} `yaml:"packages"`
	Sources      map[string]interface{} `yaml:"sources"`
	Environments map[string]interface{} `yaml:"environments"`
}

type UserConfig struct {
	Verbosity      string `yaml:"verbosity"`
	ParallelBuilds int    `yaml:"parallelBuilds"`
}
