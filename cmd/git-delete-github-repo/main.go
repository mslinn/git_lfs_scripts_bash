package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mslinn/git_lfs_scripts/internal/common"
	"github.com/mslinn/git_lfs_scripts/internal/github"
)

func main() {
	showHelp := flag.Bool("h", false, "Show help")
	flag.Parse()

	if *showHelp {
		printHelp("")
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		printHelp("Error: The name of your GitHub repository must be specified")
		os.Exit(1)
	}

	repoName := flag.Arg(0)

	// Check if gh is installed
	if err := github.CheckGHInstalled(); err != nil {
		common.PrintError("%v", err)
	}

	fmt.Printf("Deleting GitHub repository: %s\n", repoName)

	if err := github.DeleteRepo(repoName); err != nil {
		common.PrintError("%v", err)
	}

	fmt.Printf("Successfully deleted repository: %s\n", repoName)
}

func printHelp(msg string) {
	if msg != "" {
		fmt.Println(msg)
		fmt.Println()
	}

	fmt.Println("git-delete-github-repo - Delete a GitHub repository")
	fmt.Println()
	fmt.Println("Syntax: git delete-github-repo [OPTIONS] REPOSITORY_NAME")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -h  Show this help message")
	fmt.Println()
	fmt.Println("This command uses the GitHub CLI (gh) to delete a repository.")
	fmt.Println("You must have gh installed and authenticated.")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  git delete-github-repo my-test-repo")
}
