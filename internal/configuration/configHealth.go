package configuration

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/thisismeamir/hepsw/internal/utils"
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

	// Does the configuration exist?
	if _, err := os.Stat(hepswConfigPath); os.IsNotExist(err) {
		return err
	}
	newConfig := viper.New()
	newConfig.SetConfigType("yaml")
	newConfig.SetConfigName("hepsw")
	newConfig.AddConfigPath("/home/kid-a/.hepsw")

	if err := newConfig.ReadInConfig(); err != nil {
		return err
	}

	var config Configuration
	unmarshallingError := newConfig.Unmarshal(&config)
	if unmarshallingError != nil {
		return unmarshallingError
	}

	info, directoryError := config.CheckDirectories()
	if directoryError != nil {
		return errors.New(directoryError.Error() + " named " + info)
	}

	stateError := config.CheckState()
	if stateError != nil {
		return stateError
	}

	userConfigError := config.CheckUserConfigurations()
	if userConfigError != nil {
		return userConfigError
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
	if _, err := utils.CheckDirectory(config.Index); err != nil {
		return "index", err
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
