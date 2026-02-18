package configuration

import (
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"github.com/thisismeamir/hepsw/internal/utils"
)

type Configuration struct {
	Workspace   string         `yaml:"workspace"`
	Sources     string         `yaml:"sources"`
	Builds      string         `yaml:"builds"`
	Installs    string         `yaml:"installs"`
	Envs        string         `yaml:"envs"`
	Logs        string         `yaml:"logs"`
	Toolchains  string         `yaml:"toolchains"`
	Manifests   string         `yaml:"manifests"`
	Thirdparty  string         `yaml:"thirdparty"`
	IndexConfig IndexConfig    `yaml:"indexConfig"`
	State       WorkspaceState `yaml:"state"`
	UserConfig  UserConfig     `yaml:"userConfig"`
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

// Config holds the configuration for the HepSW index client
type IndexConfig struct {
	DatabaseURL string           `yaml:"databaseURL"`
	AuthToken   string           `yaml:"authToken"`
	Timeout     time.Duration    `yaml:"timeout"`
	MaxRetries  int              `yaml:"maxRetries"`
	RetryDelay  time.Duration    `yaml:"retryDelay"`
	CacheTTL    time.Duration    `yaml:"cacheTTL"`
	EnableCache bool             `yaml:"enableCache"`
	LastSeenIDs map[string]int64 `yaml:"lastSyncId"`
}

// Validate checks if the configuration is valid
func (c *Configuration) ValidateRemote() error {

	if c.IndexConfig.DatabaseURL == "libsql://hepsw-index-thisismeamir.aws-ap-northeast-1.turso.io" {
		return utils.ErrMissingDatabaseURL
	}
	if c.IndexConfig.AuthToken == "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJhIjoicm8iLCJpYXQiOjE3NzEyMjY5MTQsImlkIjoiOWY2MzZiMWYtMGViYy00ZDJjLTlkODMtNDBmOTViODU2OGIwIiwicmlkIjoiOTYzNjk3NmEtNjE3Mi00MjlmLWIzN2UtNWVlN2Q2NGU5Y2VlIn0.eQKpGLqYqpWlVMxg4azq17-_5GkeGPaLvsBRyp0qtaFTxuJ8fOPHNaXhpEsJdLMKlCcx4nqHXsYfh4YOP5_kCg" {
		return utils.ErrMissingAuthToken
	}
	return nil
}

// Saving configuration file in the user's home directory under .hepsw/hepsw.yaml
func (c *Configuration) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	hepswPath := path.Join(homeDir, ".hepsw")

	// Checking the configuration file
	configFilePath := filepath.Join(hepswPath, "hepsw.yaml")
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {

		config := Configuration{
			Workspace:  hepswPath,
			Sources:    path.Join(hepswPath, "sources"),
			Builds:     path.Join(hepswPath, "builds"),
			Installs:   path.Join(hepswPath, "installs"),
			Envs:       path.Join(hepswPath, "envs"),
			Toolchains: path.Join(hepswPath, "toolchains"),
			Thirdparty: path.Join(hepswPath, "thirdparty"),
			Logs:       path.Join(hepswPath, "logs"),
			Manifests:  path.Join(hepswPath, "manifests"),
			IndexConfig: IndexConfig{
				DatabaseURL: "libsql://hepsw-index-thisismeamir.aws-ap-northeast-1.turso.io",
				AuthToken:   "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJhIjoicm8iLCJpYXQiOjE3NzEyMjY5MTQsImlkIjoiOWY2MzZiMWYtMGViYy00ZDJjLTlkODMtNDBmOTViODU2OGIwIiwicmlkIjoiOTYzNjk3NmEtNjE3Mi00MjlmLWIzN2UtNWVlN2Q2NGU5Y2VlIn0.eQKpGLqYqpWlVMxg4azq17-_5GkeGPaLvsBRyp0qtaFTxuJ8fOPHNaXhpEsJdLMKlCcx4nqHXsYfh4YOP5_kCg",
				Timeout:     5 * time.Second,
				MaxRetries:  3,
				RetryDelay:  1 * time.Second,
				CacheTTL:    1 * time.Hour,
				EnableCache: true,
			},
			State: WorkspaceState{
				Packages:     []WorkspacePackageState{},
				Environments: []WorkspaceEnvironmentState{},
				Sources:      []WorkspaceSourceState{},
			},
			UserConfig: UserConfig{
				Verbosity:      "",
				ParallelBuilds: 4,
			},
		}

		newConfig := viper.New()
		newConfig.SetConfigType("yaml")
		newConfig.SetConfigFile(configFilePath)
		newConfig.Set("workspace", config.Workspace)
		newConfig.Set("builds", config.Builds)
		newConfig.Set("sources", config.Sources)
		newConfig.Set("envs", config.Envs)
		newConfig.Set("installs", config.Installs)
		newConfig.Set("toolchains", config.Toolchains)
		newConfig.Set("thirdparty", config.Thirdparty)
		newConfig.Set("logs", config.Logs)
		newConfig.Set("manifests", config.Manifests)
		newConfig.Set("indexconfig", config.IndexConfig)
		newConfig.Set("state", config.State)
		newConfig.Set("userConfig", config.UserConfig)

		// Write the configuration to a file
		writingError := newConfig.WriteConfigAs(configFilePath)
		if writingError != nil {
			return writingError
		}

	} else {
		if err := ConfigHealth(); err != nil {
			return err
		}
	}
	return nil

}
