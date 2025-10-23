package github

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/lithammer/dedent"
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

// CheckGHInstalled checks if the gh CLI is installed and attempts to install it if not
func CheckGHInstalled() error {
	// Check if gh is already installed
	cmd := exec.Command("gh", "--version")
	if err := cmd.Run(); err == nil {
		return nil // gh is already installed
	}

	// gh not found, attempt installation
	fmt.Println("GitHub CLI (gh) not found. Attempting to install...")

	if err := installGH(); err != nil {
		return fmt.Errorf("failed to install gh CLI: %v\nPlease install manually from: https://cli.github.com/", err)
	}

	// Verify installation succeeded
	cmd = exec.Command("gh", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gh CLI installation appeared to succeed but gh is still not available\nPlease install manually from: https://cli.github.com/")
	}

	fmt.Println("Successfully installed GitHub CLI (gh)")
	return nil
}

// installGH attempts to install gh CLI on Ubuntu and macOS
func installGH() error {
	switch runtime.GOOS {
	case "linux":
		return installGHLinux()
	case "darwin":
		return installGHMacOS()
	default:
		return fmt.Errorf("automatic installation not supported on %s", runtime.GOOS)
	}
}

// installGHLinux installs gh on Ubuntu/Debian-based systems
func installGHLinux() error {
	// Check if this is Ubuntu/Debian by looking for apt-get
	if _, err := exec.LookPath("apt-get"); err != nil {
		return fmt.Errorf("automatic installation only supported on Ubuntu/Debian (apt-get not found)")
	}

	// Check if running as root or if sudo is available
	isSudo := os.Geteuid() == 0
	if !isSudo {
		if _, err := exec.LookPath("sudo"); err != nil {
			return fmt.Errorf("sudo not available and not running as root")
		}
	}

	fmt.Println("Installing gh CLI on Ubuntu/Debian...")

	// Execute installation script as a single command
	// Based on: https://github.com/cli/cli/blob/trunk/docs/install_linux.md
	script := dedent.Dedent(`
		type -p curl >/dev/null || (sudo apt-get update && sudo apt-get install -y curl)
		curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
		sudo chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg
		echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
		sudo apt-get update
		sudo apt-get install -y gh
		`)
	if isSudo {
		script = strings.ReplaceAll(script, "sudo ", "")
	}

	cmd := exec.Command("sh", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("installation script failed: %v", err)
	}

	return nil
}

// installGHMacOS installs gh on macOS using Homebrew
func installGHMacOS() error {
	// Check if Homebrew is installed
	if _, err := exec.LookPath("brew"); err != nil {
		return fmt.Errorf("Homebrew not found. Please install Homebrew first: https://brew.sh/")
	}

	fmt.Println("Installing gh CLI using Homebrew...")

	cmd := exec.Command("brew", "install", "gh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("brew install gh failed: %v", err)
	}

	return nil
}
