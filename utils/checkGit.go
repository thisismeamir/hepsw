package utils

import "os/exec"

func CheckGit() bool {
	// check if git is available in the system
	output, err := exec.Command("git", "--version").Output()
	if err != nil {
		return false
	}
	return len(output) > 0
}
