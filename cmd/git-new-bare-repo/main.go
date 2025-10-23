package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mslinn/git_lfs_scripts/internal/common"
)

func main() {
	showHelp := flag.Bool("h", false, "Show help")
	flag.Parse()

	if *showHelp || flag.NArg() == 0 {
		printHelp("")
		os.Exit(0)
	}

	repoPath := flag.Arg(0)

	// Ensure git_access group exists
	ensureGitAccessGroup()

	// Parse repo path and name
	dir := filepath.Dir(repoPath)
	name := filepath.Base(repoPath)

	// Add .git suffix if not present
	if !strings.HasSuffix(name, ".git") {
		name = name + ".git"
	}

	fullPath := filepath.Join(dir, name)

	// Check if repo already exists
	if _, err := os.Stat(fullPath); err == nil {
		printHelp(fmt.Sprintf("Error: '%s' already exists.", fullPath))
		os.Exit(1)
	}

	// Create parent directory if needed
	if err := os.MkdirAll(dir, 0755); err != nil {
		common.PrintError("Failed to create parent directory: %v", err)
	}

	// Create the bare repository directory with SGID
	fmt.Printf("Creating bare repository at %s\n", fullPath)

	if err := os.MkdirAll(fullPath, 0775); err != nil {
		common.PrintError("Failed to create repository directory: %v", err)
	}

	// Change to the parent directory
	if err := os.Chdir(dir); err != nil {
		common.PrintError("Failed to change to directory %s: %v", dir, err)
	}

	// Set group ownership to git_access (requires sudo on Linux)
	// This may fail on systems without sudo or git_access group
	cmd := exec.Command("sudo", "chgrp", "git_access", name)
	_ = cmd.Run() // Ignore error if sudo/chgrp fails

	// Initialize bare repository with shared permissions
	fmt.Println("Initializing bare repository...")
	if err := initBareRepo(fullPath); err != nil {
		common.PrintError("Failed to initialize bare repository: %v", err)
	}

	// Configure the repository
	if err := os.Chdir(fullPath); err != nil {
		common.PrintError("Failed to change to repository directory: %v", err)
	}

	fmt.Println("Configuring repository...")
	if err := configureRepo(); err != nil {
		common.PrintError("Failed to configure repository: %v", err)
	}

	fmt.Printf("Successfully created bare repository at %s\n", fullPath)
}

func printHelp(msg string) {
	if msg != "" {
		fmt.Println(msg)
		fmt.Println()
	}

	fmt.Println("git-new-bare-repo - Create a new bare Git repository")
	fmt.Println()
	fmt.Println("Syntax: git new-bare-repo [OPTIONS] /path/to/new/repo.git")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -h  Show this help message")
	fmt.Println()
	fmt.Println("Normally this script would be run on a Git server, because")
	fmt.Println("that is where bare Git repositories normally live.")
	fmt.Println()
	fmt.Println("A new Git repository will be created in /path/to/new/repo.git,")
	fmt.Println("which should not already exist.")
	fmt.Println()
	fmt.Println("The SGID permission for the new Git repository will be set for group git_access,")
	fmt.Println("which is created if it does not exist.")
	fmt.Println()
	fmt.Println("The parent directory (/path/to/new/) will be created if it does not already exist.")
	fmt.Println("The name of the repo must not contain spaces.")
	fmt.Println("If the specified name does not end with a .git suffix, the suffix is appended.")
	fmt.Println()
	fmt.Println("Git configuration parameter 'receive.denyCurrentBranch' is set to ignore.")
}

func ensureGitAccessGroup() {
	// Check if git_access group exists, create if needed
	cmd := exec.Command("getent", "group", "git_access")
	if err := cmd.Run(); err != nil {
		// Group doesn't exist, try to create it
		createCmd := exec.Command("sudo", "groupadd", "git_access")
		_ = createCmd.Run() // Ignore error if this fails
	}
}

func initBareRepo(path string) error {
	cmd := exec.Command("git", "init", "--bare", "--shared=everybody", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func configureRepo() error {
	cmd := exec.Command("git", "config", "receive.denyCurrentBranch", "ignore")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
