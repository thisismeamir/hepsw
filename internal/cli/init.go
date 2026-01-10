package cli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/utils"
)

// Initialization command that sets up a new HepSW workspace, configures environment variables, and creates necessary directories and files.

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize HepSW workspace",
	Long: `The first step in using HepSW is to initialize a workspace. 
This workspace will contain all the necessary directories and 
configuration files needed to manage your software stack. to 
initialize the workspace use:

	hepsw init
`,
	RunE: runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	PrintSection("Initializing HepSW...")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		PrintError("Error getting home directory got" + " " + err.Error())
		return err
	}

	// Creating ~/.hepsw directory
	hepswPath := homeDir + "/.hepsw"
	err = utils.CreateDirectory(hepswPath)
	if err != nil {
		PrintError("Error creating HepSW directory " + err.Error())
		return err
	} else {
		PrintSuccess("HepSW directory is ready!")
	}

	// Creating subdirectories
	subdirectories := []string{"toolchains", "sources", "builds", "installs", "envs", "logs", "third-party"}

	for _, item := range subdirectories {
		err = utils.CreateDirectory(hepswPath + "/" + item)
		if err != nil {
			PrintError("Error creating HepSW directory " + err.Error())
			return err
		} else {
			PrintSuccess(hepswPath + "/" + item + " directory is ready!")
		}
	}

	// Syncing local index with the remote repository.
	_ = Sync()
	_ = ConfigInit()
	return nil
}
