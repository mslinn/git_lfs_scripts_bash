package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/lithammer/dedent"
	"github.com/mslinn/git_lfs_scripts/internal/common"
	"github.com/mslinn/git_lfs_scripts/internal/lfsfiles"
	flag "github.com/spf13/pflag"
)

func main() {
	var bothCases, dryRun, everywhere, showHelp bool

	flag.BoolVarP(&bothCases, "case", "c", false, "Expand pattern to upper and lower case")
	flag.BoolVarP(&dryRun, "dry-run", "d", false, "Dry run")
	flag.BoolVarP(&everywhere, "everywhere", "e", false, "Apply pattern everywhere")
	flag.BoolVarP(&showHelp, "help", "h", false, "Show help")
	flag.Parse()

	if showHelp {
		printHelp()
		os.Exit(0)
	}

	patterns := flag.Args()
	if len(patterns) == 0 {
		printHelp()
		os.Exit(1)
	}

	// Check if we're in a git repository
	if err := common.CheckGitRepo(); err != nil {
		common.PrintError("%v", err)
	}

	// Check if git-lfs is installed
	if err := common.CheckLFSInstalled(); err != nil {
		common.PrintError("%v", err)
	}

	// Check if LFS is initialized in this repo
	if err := common.CheckLFSInitialized(); err != nil {
		common.PrintError("%v", err)
	}

	opts := lfsfiles.Options{
		BothCases:  bothCases,
		DryRun:     dryRun,
		Everywhere: everywhere,
		Command:    "git lfs untrack",
	}

	// If dry run, just show what would be done
	if dryRun {
		for _, pattern := range patterns {
			expanded := lfsfiles.ExpandPattern(pattern, opts)
			fmt.Printf("DRY RUN: git lfs untrack %s\n", strings.Join(expanded, " "))
		}
		fmt.Println("DRY RUN: git add --renormalize .")
		fmt.Printf("DRY RUN: git commit -m \"Restore patterns to Git from Git LFS\"\n")
		fmt.Println("DRY RUN: git push")
		os.Exit(0)
	}

	// Untrack patterns from LFS
	for _, pattern := range patterns {
		expanded := lfsfiles.ExpandPattern(pattern, opts)
		args := append([]string{"lfs", "untrack"}, expanded...)

		cmd := exec.Command("git", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			common.PrintError("Failed to untrack pattern %s: %v", pattern, err)
		}
	}

	// Renormalize and commit
	fmt.Println("Renormalizing files...")
	if err := runGitCommand("add", "--renormalize", "."); err != nil {
		common.PrintError("Failed to renormalize: %v", err)
	}

	commitMsg := fmt.Sprintf("Restore patterns to Git from Git LFS")
	fmt.Printf("Committing changes...\n")
	if err := runGitCommand("commit", "-m", commitMsg); err != nil {
		// It's ok if there's nothing to commit
		fmt.Println("No changes to commit")
	}

	fmt.Println("Pushing changes...")
	if err := runGitCommand("push"); err != nil {
		common.PrintError("Failed to push: %v", err)
	}

	fmt.Println("Unmigration complete!")
}

func printHelp() {
	fmt.Print(dedent.Dedent(`
		git-unmigrate - Move matching files from Git LFS back to Git

		USAGE:
		  git unmigrate [OPTIONS] PATTERN ...

		OPTIONS:
		  -c  Expand pattern to upper and lower case, helpful for media files
		  -d  Dry run (display filename patterns that would be affected)
		  -e  Apply the pattern everywhere (all directories in the Git repository)
		  -h  Show this help message

		DESCRIPTION:
		  This command reverses 'git lfs migrate import' by moving files back to regular
		  Git tracking. By default, only files in the current directory matching the
		  specified patterns are processed.

		  This process does NOT rewrite Git history, so other Git users will not need
		  to re-clone the repository after this process concludes.

		  Note: This process might take a long time if you have many large files to
		  unmigrate back to Git.

		REQUIREMENTS:
		  - Git repository
		  - Git LFS installed and configured

		EXAMPLES:
		  # Dry run - see what would happen
		  git unmigrate -d zip
		  # Output: DRY RUN: git lfs untrack *.zip

		  # Unmigrate multiple patterns
		  git unmigrate -d pdf zip
		  # Output: DRY RUN: git lfs untrack *.pdf
		  #         DRY RUN: git lfs untrack *.zip

		  # Unmigrate with case variations
		  git unmigrate -dc mp3
		  # Output: DRY RUN: git lfs untrack *.mp3 *.MP3

		  # Apply everywhere in repository
		  git unmigrate -de zip
		  # Output: DRY RUN: git lfs untrack *.zip **/*.zip

		  # Combined options
		  git unmigrate -dce mp3
		  # Output: DRY RUN: git lfs untrack *.mp3 *.MP3 **/*.mp3 **/*.MP3

		  # Actually unmigrate (remove -d flag)
		  git unmigrate zip

		SEE ALSO:
		  git-ls-files, git-lfs-track, git-lfs-untrack
	`))
}

func runGitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
