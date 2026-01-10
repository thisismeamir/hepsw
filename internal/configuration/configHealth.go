package configuration

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/thisismeamir/hepsw/internal/utils"
)

type ConfigurationHealthType interface {
}

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
	viper.SetConfigName("hepsw")
	viper.AddConfigPath(hepswConfigPath)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	var config Configuration
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	directoryError := config.CheckDirectories()
	if directoryError != nil {
		return directoryError
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

func (config *Configuration) CheckDirectories() error {
	// Checking that all the directories in the workspace exist
	if _, err := utils.CheckDirectory(config.Workspace); err != nil {
		return err
	}
	if _, err := utils.CheckDirectory(config.SourcesDir); err != nil {
		return err
	}
	if _, err := utils.CheckDirectory(config.BuildsDir); err != nil {
		return err
	}
	if _, err := utils.CheckDirectory(config.InstallsDir); err != nil {
		return err
	}
	if _, err := utils.CheckDirectory(config.EnvsDir); err != nil {
		return err
	}
	if _, err := utils.CheckDirectory(config.LogsDir); err != nil {
		return err
	}
	if _, err := utils.CheckDirectory(config.ToolchainsDir); err != nil {
		return err
	}
	if _, err := utils.CheckDirectory(config.indexDir); err != nil {
		return err
	}
	return nil
}

func (config *Configuration) CheckState() error {
	// TODO: State should be checked when manifest is integrated.
	return nil
}

func (config *Configuration) CheckUserConfigurations() error {
	// TODO: User Configuration is not anything important right now.
	return nil
}
