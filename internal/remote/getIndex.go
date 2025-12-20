package remote

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

// https://github.com/thisismeamir/hepsw-package-index.git modules.
// These modules help HepSW to fetch package information, versions, and dependencies from a centralized repository.

const PackageIndexRepoURL = "https://github.com/thisismeamir/hepsw-package-index.git"

// RepoExists checks if a git repository exists at the given path
func RepoExists(repoDir string) bool {
	_, err := os.Stat(repoDir)
	return !os.IsNotExist(err)
}

// OpenRepo opens an existing git repository
func OpenRepo(repoDir string) (*git.Repository, error) {
	return git.PlainOpen(repoDir)
}

// CloneRepo clones a repository from the given URL to the specified directory
func CloneRepo(repoDir, branch string) error {
	branchRef := plumbing.NewBranchReferenceName(branch)

	_, err := git.PlainClone(repoDir, false, &git.CloneOptions{
		URL:           PackageIndexRepoURL,
		ReferenceName: branchRef,
		SingleBranch:  true,
	})
	return err
}

// HasLocalChanges checks if the repository has any uncommitted changes
func HasLocalChanges(repo *git.Repository) (bool, error) {
	worktree, err := repo.Worktree()
	if err != nil {
		return false, err
	}

	status, err := worktree.Status()
	if err != nil {
		return false, err
	}

	return !status.IsClean(), nil
}

// ResetLocalChanges performs a hard reset to HEAD and cleans untracked files
func ResetLocalChanges(repo *git.Repository) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Hard reset to HEAD
	err = worktree.Reset(&git.ResetOptions{
		Mode: git.HardReset,
	})
	if err != nil {
		return err
	}

	// Clean untracked files
	return worktree.Clean(&git.CleanOptions{
		Dir: true,
	})
}

// FetchRemote fetches the latest changes from the remote repository
func FetchRemote(repo *git.Repository, branch string) error {
	branchRef := plumbing.NewBranchReferenceName(branch)
	remote, err := repo.Remote("origin")
	if err != nil {
		return err
	}

	return remote.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec(fmt.Sprintf("+%s:%s", branchRef, branchRef)),
		},
	})
}

// HasRemoteUpdates checks if there are updates available from the remote
func HasRemoteUpdates(repo *git.Repository, branch string) (bool, error) {
	branchRef := plumbing.NewBranchReferenceName(branch)

	// Get HEAD reference
	ref, err := repo.Head()
	if err != nil {
		return false, err
	}

	// Get the remote reference
	remoteRef, err := repo.Reference(branchRef, true)
	if err != nil {
		return false, err
	}

	return ref.Hash() != remoteRef.Hash(), nil
}

// PullChanges pulls the latest changes from the remote repository
func PullChanges(repo *git.Repository, branch string) error {
	branchRef := plumbing.NewBranchReferenceName(branch)
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	return worktree.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: branchRef,
		Force:         true,
	})
}
