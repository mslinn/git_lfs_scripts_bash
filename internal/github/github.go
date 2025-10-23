package github

import (
	"fmt"
	"os/exec"
)

// DeleteRepo deletes a GitHub repository using the gh CLI
func DeleteRepo(repoName string) error {
	cmd := exec.Command("gh", "repo", "delete", repoName, "--yes")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to delete repository %s: %v\nOutput: %s", repoName, err, string(output))
	}

	return nil
}

// CheckGHInstalled checks if the gh CLI is installed
func CheckGHInstalled() error {
	cmd := exec.Command("gh", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gh CLI is not installed or not in PATH. Install from: https://cli.github.com/")
	}
	return nil
}
