package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lithammer/dedent"
	"github.com/mslinn/git_lfs_scripts/internal/common"
	flag "github.com/spf13/pflag"
)

func main() {
	showHelp := flag.BoolP("help", "h", false, "Show help")
	flag.Parse()

	if *showHelp || flag.NArg() == 0 {
		printHelp("")
		os.Exit(0)
	}

	repoPath := flag.Arg(0)

	// Validate input
	if repoPath == "." || repoPath == ".." || repoPath == "/" {
		printHelp(fmt.Sprintf("Error: Invalid repository path '%s'.\nPlease provide a specific repository name or path.", repoPath))
		os.Exit(1)
	}

	// Check prerequisites
	checkPrerequisites()

	// Ensure git_access group exists
	ensureGitAccessGroup()

	// Parse repo path and name
	// Clean the path first to handle relative paths properly
	cleanPath := filepath.Clean(repoPath)
	dir := filepath.Dir(cleanPath)
	name := filepath.Base(cleanPath)

	// Add .git suffix if not present
	if !strings.HasSuffix(name, ".git") {
		name = name + ".git"
	}

	fullPath := filepath.Join(dir, name)

	// Convert to absolute path for clarity
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		common.PrintError("Failed to resolve absolute path: %v", err)
	}
	fullPath = absPath

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

	fmt.Print(dedent.Dedent(`
		git-new-bare-repo - Create a new bare Git repository

		USAGE:
		  git new-bare-repo [OPTIONS] /path/to/new/repo.git

		OPTIONS:
		  -h  Show this help message

		DESCRIPTION:
		  Creates a new bare Git repository, typically run on a Git server where bare
		  repositories normally live.

		  The new repository will be created at the specified path (which must not
		  already exist). The SGID permission will be set for group git_access, which
		  is created if it does not exist.

		  Features:
		    - Parent directories are created automatically if needed
		    - .git suffix is appended if not specified
		    - Shared repository permissions (--shared=everybody)
		    - Sets receive.denyCurrentBranch to ignore

		  Note: Repository names must not contain spaces.

		REQUIREMENTS:
		  - Git
		  - sudo (for group management operations)
		  - getent (for checking group existence)
		  - groupadd (for creating git_access group)
		  - chgrp (for setting group ownership)

		EXAMPLES:
		  # Create a repository (adds .git automatically)
		  git new-bare-repo /srv/git/myproject

		  # Create with explicit .git suffix
		  git new-bare-repo /srv/git/myproject.git

		  # Create in a nested path (parent dirs created automatically)
		  git new-bare-repo /srv/git/team/project.git
	`))
}

func checkPrerequisites() {
	var missing []string

	// Check git
	if _, err := exec.LookPath("git"); err != nil {
		missing = append(missing, "git (install from: https://git-scm.com/)")
	}

	// Check sudo
	if _, err := exec.LookPath("sudo"); err != nil {
		missing = append(missing, "sudo (required for group management)")
	}

	// Check getent
	if _, err := exec.LookPath("getent"); err != nil {
		missing = append(missing, "getent (usually part of glibc-common)")
	}

	// Check groupadd
	if _, err := exec.LookPath("groupadd"); err != nil {
		missing = append(missing, "groupadd (usually part of shadow-utils)")
	}

	// Check chgrp
	if _, err := exec.LookPath("chgrp"); err != nil {
		missing = append(missing, "chgrp (usually part of coreutils)")
	}

	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "Error: Missing required commands:\n")
		for _, cmd := range missing {
			fmt.Fprintf(os.Stderr, "  ✗ %s\n", cmd)
		}
		fmt.Fprintf(os.Stderr, "\nPlease install missing dependencies before running git-new-bare-repo.\n")
		os.Exit(1)
	}
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
