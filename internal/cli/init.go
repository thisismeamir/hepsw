package cli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/remote"
	"github.com/thisismeamir/hepsw/utils"
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
	PrintSection("Initializing HepSW...")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		PrintError("Error getting home directory got" + " " + err.Error())
		os.Exit(1)
	}

	// Creating ~/.hepsw directory
	hepswPath := homeDir + "/.hepsw"
	err = utils.CreateDirectory(hepswPath)
	if err != nil {
		PrintError("Error creating HepSW directory " + err.Error())
		os.Exit(1)
	} else {
		PrintSuccess("HepSW directory is ready!")
	}

	// Creating subdirectories
	subdirectories := []string{"logs", "builds", "sources", "installs", "envs"}

	for _, item := range subdirectories {
		err = utils.CreateDirectory(hepswPath + "/" + item)
		if err != nil {
			PrintError("Error creating HepSW directory " + err.Error())
			os.Exit(1)
		} else {
			PrintSuccess(hepswPath + "/" + item + " directory is ready!")
		}
	}

	// Fetching and initializing index repo
	repoDir := hepswPath + "/index"
	// Check if the index repo exists
	repoExists := remote.RepoExists(repoDir)
	if repoExists {
		repo, err := remote.OpenRepo(repoDir)
		if err != nil {
			PrintError("Error opening repo " + err.Error())
			os.Exit(1)
		}
		repoHasChanges, err := remote.HasLocalChanges(repo)
		if err != nil {
			PrintError("Error checking if repo has changes " + err.Error())
		}
		if repoHasChanges {
			PrintWarning("You should not edit the content of index repo, resetting local changes...")
			err := remote.ResetLocalChanges(repo)
			if err != nil {
				PrintError("Error resetting local changes " + err.Error())
			}
		}
		repoHasUpdate, err := remote.HasRemoteUpdates(repo, "master")
		if repoHasUpdate {
			PrintInfo("Index repo has updates.")
			err := remote.FetchRemote(repo, "master")
			if err != nil {
				PrintError("Error fetching remote " + err.Error())
				PrintWarning("Might not be update repository for better indexing.")
			}
			err = remote.PullChanges(repo, "master")
			if err != nil {
				PrintError("Error pulling changes " + err.Error())
				PrintWarning("Might not be update repository for better indexing.")
			}
			PrintInfo("Index repo has been updated.")
		}
	} else {
		PrintWarning("HepSW directory does not contain index repository cloning...")
		err := remote.CloneRepo(repoDir, "master")
		if err != nil {
			PrintError("Error cloning repo " + err.Error())
			os.Exit(1)
		}
	}
}
