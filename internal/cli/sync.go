package cli

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/remote"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronizes changes in index repository with the ones on remote.",
	Long: `Synchronizes changes in index repository with the ones on remote.
if a change has been made to this repository the program would automatically remove it, ensuring that the repository stays 
relevant to the remote one.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Sync()
	},
}

func Sync() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	hepswPath := filepath.Join(homeDir, ".hepsw")
	// Fetching and initializing index repo
	repoDir := hepswPath + "/index"
	// Check if the index repo exists
	repoExists := remote.RepoExists(repoDir)
	if repoExists {
		repo, err := remote.OpenRepo(repoDir)
		if err != nil {
			PrintError("Error opening repo " + err.Error())
			return err
		}
		repoHasChanges, err := remote.HasLocalChanges(repo)
		if err != nil {
			PrintError("Error checking if repo has changes " + err.Error())
			return err
		}
		if repoHasChanges {
			PrintWarning("You should not edit the content of index repo, resetting local changes...")
			err := remote.ResetLocalChanges(repo)
			if err != nil {
				PrintError("Error resetting local changes " + err.Error())
				return err
			}
		}
		repoHasUpdate, err := remote.HasRemoteUpdates(repo, "master")
		if repoHasUpdate {
			PrintInfo("Index repo has updates.")
			err := remote.FetchRemote(repo, "master")
			if err != nil {
				PrintError("Error fetching remote " + err.Error())
				PrintWarning("Might not be update repository for better indexing.")
				return err
			}
			err = remote.PullChanges(repo, "master")
			if err != nil {
				PrintError("Error pulling changes " + err.Error())
				PrintWarning("Might not be update repository for better indexing.")
				return err
			}
			PrintInfo("Index repo has been updated.")
		}
	} else {
		PrintWarning("HepSW directory does not contain index repository cloning...")
		err := remote.CloneRepo(repoDir, "master")
		if err != nil {
			PrintError("Error cloning repo " + err.Error())
			return err
		}
	}

	return nil
}
