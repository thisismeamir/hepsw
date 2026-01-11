package cli

import (
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thisismeamir/hepsw/internal/configuration"
)

var checkConfig = &cobra.Command{
	Use:   "check-config",
	Short: "Checks if the configuration is available, otherwise creates it.",
	Long:  `Checks if the configuration is available, otherwise creates it.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ConfigInit()
	},
}

func ConfigInit() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	hepswPath := path.Join(homeDir, ".hepsw")

	// Checking the configuration file
	configFilePath := filepath.Join(hepswPath, "hepsw.yaml")
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {

		config := configuration.Configuration{
			Workspace:  hepswPath,
			Sources:    path.Join(hepswPath, "sources"),
			Builds:     path.Join(hepswPath, "builds"),
			Installs:   path.Join(hepswPath, "installs"),
			Envs:       path.Join(hepswPath, "envs"),
			Toolchains: path.Join(hepswPath, "toolchains"),
			Thirdparty: path.Join(hepswPath, "thirdparty"),
			Logs:       path.Join(hepswPath, "logs"),
			Index:      path.Join(hepswPath, "index"),
			State: configuration.WorkspaceState{
				Packages:     []configuration.WorkspacePackageState{},
				Environments: []configuration.WorkspaceEnvironmentState{},
				Sources:      []configuration.WorkspaceSourceState{},
			},
			UserConfig: configuration.UserConfig{
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
		newConfig.Set("index", config.Index)
		newConfig.Set("state", config.State)
		newConfig.Set("userConfig", config.UserConfig)

		// Write the configuration to a file
		writingError := newConfig.WriteConfigAs(configFilePath)
		if writingError != nil {
			PrintError("Error writing config file: " + writingError.Error())
			return writingError
		}

	} else {
		if err := configuration.ConfigHealth(); err != nil {
			PrintError("Configuration is not healthy: " + err.Error())
			return err
		}
		PrintSuccess("Configuration has been loaded successfully.")
	}
	return nil
}
