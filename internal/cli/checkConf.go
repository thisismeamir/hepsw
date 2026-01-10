package cli

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var resetCmd = &cobra.Command{
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
		// Creating the HepSW configuration file ~/.hepsw/hepsw.yaml
		config := viper.New()
		config.SetConfigType("yaml")
		config.SetConfigName("hepsw")
		config.AddConfigPath(hepswPath)
		config.Set("workspace", hepswPath)
		config.Set("sourcesDir", path.Join(hepswPath, "sources"))
		config.Set("buildsDir", path.Join(hepswPath, "builds"))
		config.Set("installsDir", path.Join(hepswPath, "installs"))
		config.Set("envsDir", path.Join(hepswPath, "envs"))
		config.Set("logsDir", path.Join(hepswPath, "logs"))
		config.Set("thirdPartyDir", path.Join(hepswPath, "third_party"))
		config.Set("toolchainsDir", path.Join(hepswPath, "toolchains"))
		config.Set("indexDir", path.Join(hepswPath, "index"))

		state := map[string]interface{}{
			"packages":     map[string]interface{}{},
			"sources":      map[string]interface{}{},
			"environments": map[string]interface{}{},
		}

		userConfig := map[string]interface{}{
			"verbosity":   map[string]interface{}{},
			"parallelism": map[string]interface{}{},
		}

		config.Set("state", state)
		config.Set("userConfig", userConfig)

		// Write the configuration to a file
		err := config.WriteConfigAs(configFilePath)
		if err != nil {
			PrintError("Error writing config file: " + err.Error())
			return err
		}

		// Verify that the configuration was written correctly
		if err := config.ReadInConfig(); err == nil {
			PrintSuccess("Configuration loaded successfully:")
			fmt.Println(config.AllSettings())
		} else {
			PrintError("Error loading configuration:" + err.Error())
			return err
		}
	} else {
		PrintWarning("Configuration is loaded successfully.")
	}
	return nil
}
