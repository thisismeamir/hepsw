package configuration

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/thisismeamir/hepsw/internal/utils"
	"gopkg.in/yaml.v3"
)

// ConfigHealth
// Checks the configuration file health by checking that the directories exist, as well as
// checking the state of the workspace, packages installed, environments and other important
// stuff.
func ConfigHealth() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	hepswConfigPath := filepath.Join(homeDir, ".hepsw", "hepsw.yaml")
	data, err2 := os.ReadFile(hepswConfigPath)
	if err2 != nil {
		return err2
	}

	var config Configuration
	if yamlError := yaml.Unmarshal(data, &config); yamlError != nil {
		return yamlError
	}

	// Check State
	if err := config.CheckState(); err != nil {
		return errors.New("State is not healthy: " + err.Error())
	}
	// Check Directories
	if str, err := config.CheckDirectories(); err != nil {
		return errors.New("Directory " + str + " is not healthy: " + err.Error())
	}

	// Check User Configurations
	if err := config.CheckUserConfigurations(); err != nil {
		return errors.New("User configurations are not healthy: " + err.Error())
	}
	return nil
}

func (config *Configuration) CheckDirectories() (string, error) {
	// Checking that all the directories in the workspace exist
	if _, err := utils.CheckDirectory(config.Workspace); err != nil {
		return "workspace", err
	}
	if _, err := utils.CheckDirectory(config.Sources); err != nil {
		return "sources", err
	}
	if _, err := utils.CheckDirectory(config.Builds); err != nil {
		return "builds", err
	}
	if _, err := utils.CheckDirectory(config.Installs); err != nil {
		return "installs", err
	}
	if _, err := utils.CheckDirectory(config.Envs); err != nil {
		return "envs", err
	}
	if _, err := utils.CheckDirectory(config.Logs); err != nil {
		return "logs", err
	}
	if _, err := utils.CheckDirectory(config.Toolchains); err != nil {
		return "toolchains", err
	}
	if _, err := utils.CheckDirectory(config.Manifests); err != nil {
		return "manifests", err
	}
	if _, err := utils.CheckDirectory(config.Thirdparty); err != nil {
		return "thirdparty", err
	}
	return "", nil
}

func (config *Configuration) CheckState() error {
	// TODO: State should be checked when manifest is integrated.
	return nil
}

func (config *Configuration) CheckUserConfigurations() error {
	// TODO: User Configuration is not anything important right now.
	return nil
}
