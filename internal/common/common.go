package common

import (
	"fmt"
	"os"
	"os/exec"
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
