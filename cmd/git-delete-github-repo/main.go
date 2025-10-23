package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lithammer/dedent"
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

	fmt.Print(dedent.Dedent(`
		git-delete-github-repo - Delete a GitHub repository

		SYNTAX:
		  git delete-github-repo [OPTIONS] REPOSITORY_NAME

		OPTIONS:
		  -h  Show this help message

		DESCRIPTION:
		  This command uses the GitHub CLI (gh) to delete a repository.

		  If gh is not installed, it will attempt automatic installation on:
		    - Ubuntu/Debian (using apt-get)
		    - macOS (using Homebrew)

		  You must have gh authenticated (run 'gh auth login' after installation).

		EXAMPLE:
		  git delete-github-repo my-test-repo
	`))
}
