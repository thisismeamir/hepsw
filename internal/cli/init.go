package cli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Initialization command that sets up a new HepSW workspace, configures environment variables, and creates necessary directories and files.

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new HepSW workspace",
	Long: `Initialize a new HepSW workspace by creating necessary directories,
	configuration files, and setting up environment variables.`,
	Run: runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	// Assert exactly one argument is provided (workspace directory)
	if len(args) != 1 {
		PrintError("For initialization, please provide exactly one argument: the path to the workspace directory.\nExample: hepsw init /path/to/workspace")
		os.Exit(1)
	}

	PrintSection("HepSW Workspace initialization")
	// Implementation of the initialization logic goes here
	workSpaceDir := args[0]
	// assert workspace directory exists or create it
	if _, err := os.Stat(workSpaceDir); os.IsNotExist(err) {
		PrintInfo("Creating workspace directory: " + workSpaceDir)
		err := os.MkdirAll(workSpaceDir, 0755)
		if err != nil {
			panic(err)
		}
	} else {
		PrintWarning("Workspace directory already exists: " + workSpaceDir)
	}
	PrintSection("Looking for subdirectories...")
	// Creating necessary subdirectories
	subDirs := []string{"packages", "build", "install", "logs"}
	for _, dir := range subDirs {
		fullPath := workSpaceDir + "/" + dir
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			PrintInfo("Creating subdirectory: " + fullPath)
			err := os.MkdirAll(fullPath, 0755)
			if err != nil {
				PrintError("Something went wrong while creating subdirectory: " + fullPath)
				panic(err)
			}
		} else {
			PrintWarning("Subdirectory already exists: " + fullPath)
		}
	}

	// create config file
	configFilePath := workSpaceDir + "/hepsw.yaml"
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		PrintInfo("Creating config file: " + configFilePath)
		file, err := os.Create(configFilePath)
		if err != nil {
			panic(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}(file)

		PrintInfo("Writing defaults in config file: " + configFilePath)
		// Adding default content to config file
		configuration := viper.New()
		configuration.SetConfigType("yaml")
		configuration.Set("workspace", workSpaceDir)
		configuration.Set("packages", workSpaceDir+"/packages")
		configuration.Set("build", workSpaceDir+"/build")
		configuration.Set("installs", workSpaceDir+"/install")
		configuration.Set("logs", workSpaceDir+"/logs")

		err2 := configuration.WriteConfigAs(configFilePath)
		if err2 != nil {
			PrintError("Error writing config file: " + configFilePath)
			panic(err2)
		}

	} else {
		PrintWarning("Config file already exists: " + configFilePath)
		PrintInfo("Opening existing config file: " + configFilePath)
		configuration := viper.New()
		configuration.SetConfigFile(configFilePath)
		err := configuration.ReadInConfig()
		if err != nil {
			PrintError("Error reading config file: " + configFilePath)
			panic(err)
		}
		PrintInfo("Config file loaded successfully: " + configFilePath)
		PrintSection("Configuration: ")
		settings := configuration.AllSettings()
		for key, value := range settings {
			PrintInfo(key + ": " + value.(string))
		}
	}

}
