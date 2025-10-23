package common

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Version of the git_lfs_scripts suite
const Version = "1.0.0"

// ExecGitCommand executes a git command and returns the combined output
func ExecGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// CheckGitRepo verifies we're inside a git repository
func CheckGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("not a git repository (or any of the parent directories)")
	}
	return nil
}

// PrintError prints an error message to stderr and exits
func PrintError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}

// CheckLFSInstalled verifies Git LFS is installed
func CheckLFSInstalled() error {
	cmd := exec.Command("git", "lfs", "version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Git LFS is not installed or not available.\nInstall from: https://git-lfs.com/")
	}
	return nil
}

// CheckLFSInitialized verifies Git LFS is configured in the repository
func CheckLFSInitialized() error {
	// Check if .gitattributes exists and has LFS patterns
	file, err := os.Open(".gitattributes")
	if os.IsNotExist(err) {
		return fmt.Errorf("Git LFS is not configured for this repository.\nNo .gitattributes file found.\n\nLearn about Git LFS at:\n  https://www.mslinn.com/git/5100-git-lfs-overview.html")
	}
	if err != nil {
		return fmt.Errorf("error reading .gitattributes: %v", err)
	}
	defer file.Close()

	// Check if .gitattributes contains any LFS patterns
	scanner := bufio.NewScanner(file)
	hasLFSPattern := false
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "filter=lfs") {
			hasLFSPattern = true
			break
		}
	}

	if !hasLFSPattern {
		return fmt.Errorf("Git LFS is not configured for this repository.\nNo LFS tracked patterns found in .gitattributes.\n\nLearn about Git LFS at:\n  https://www.mslinn.com/git/5100-git-lfs-overview.html")
	}

	return nil
}
