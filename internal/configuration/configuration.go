package configuration

type Configuration struct {
	Workspace  string         `yaml:"workspace"`
	Sources    string         `yaml:"sources"`
	Builds     string         `yaml:"builds"`
	Installs   string         `yaml:"installs"`
	Envs       string         `yaml:"envs"`
	Logs       string         `yaml:"logs"`
	Toolchains string         `yaml:"toolchains"`
	Index      string         `yaml:"index"`
	Thirdparty string         `yaml:"thirdparty"`
	LastSyncId int64          `yaml:"lastSyncId"`
	State      WorkspaceState `yaml:"state"`
	UserConfig UserConfig     `yaml:"userConfig"`
}

type WorkspaceState struct {
	Packages     []WorkspacePackageState     `yaml:"packages"`
	Sources      []WorkspaceSourceState      `yaml:"sources"`
	Environments []WorkspaceEnvironmentState `yaml:"environments"`
}

type WorkspacePackageState struct {
	PackageId   string   `yaml:"packageId"`
	Name        string   `yaml:"name"`
	Path        string   `yaml:"path"`
	Version     string   `yaml:"version"`
	BuildTime   string   `yaml:"buildTime"`
	InstallTime string   `yaml:"installTime"`
	IsUsedBy    []string `yaml:"isUsedBy"`
	IsUsing     []string `yaml:"isUsing"`
}

type WorkspaceSourceState struct {
	SourceId string   `yaml:"sourceId"`
	Name     string   `yaml:"name"`
	Path     string   `yaml:"path"`
	Version  string   `yaml:"version"`
	IsUsedBy []string `yaml:"isUsedBy"`
	IsUsing  []string `yaml:"isUsing"`
}

type WorkspaceEnvironmentState struct {
	EnvironmentId string   `yaml:"environmentId"`
	Name          string   `yaml:"name"`
	Description   string   `yaml:"description"`
	ScriptPath    string   `yaml:"scriptPath"`
	Packages      []string `yaml:"packages"`
}

type UserConfig struct {
	Verbosity      string `yaml:"verbosity"`
	ParallelBuilds int    `yaml:"parallelBuilds"`
}
