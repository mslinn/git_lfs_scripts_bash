package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mslinn/git_lfs_scripts/internal/common"
	"github.com/mslinn/git_lfs_scripts/internal/lfsfiles"
)

func main() {
	var bothCases, dryRun, everywhere, showHelp bool

	flag.BoolVar(&bothCases, "c", false, "Expand pattern to upper and lower case")
	flag.BoolVar(&dryRun, "d", false, "Dry run")
	flag.BoolVar(&everywhere, "e", false, "Apply pattern everywhere")
	flag.BoolVar(&showHelp, "h", false, "Show help")
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
	fmt.Println("git-unmigrate - Move matching files from Git LFS to Git")
	fmt.Println()
	fmt.Println("By default, only files in the current directory that match one of the")
	fmt.Println("specified filetypes are processed.")
	fmt.Println()
	fmt.Println("This process does not rewrite Git history, so other Git users will not need")
	fmt.Println("to re-clone the repository after this process concludes.")
	fmt.Println()
	fmt.Println("This process might take a long time if you have a lot of large files to")
	fmt.Println("unmigrate back to Git.")
	fmt.Println()
	fmt.Println("Syntax:")
	fmt.Println("  git unmigrate [OPTIONS] PATTERN ...")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -c  Expand pattern to upper and lower case, helpful for media files")
	fmt.Println("  -d  Dry run (display filename patterns that would be affected)")
	fmt.Println("  -e  Apply the pattern everywhere (all directories in the Git repository)")
	fmt.Println("  -h  Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  $ git unmigrate -d zip")
	fmt.Println("  DRY RUN: git lfs untrack *.zip")
	fmt.Println()
	fmt.Println("  $ git unmigrate -d pdf zip")
	fmt.Println("  DRY RUN: git lfs untrack *.pdf")
	fmt.Println("  DRY RUN: git lfs untrack *.zip")
	fmt.Println()
	fmt.Println("  $ git unmigrate -dc mp3")
	fmt.Println("  DRY RUN: git lfs untrack *.mp3 *.MP3")
	fmt.Println()
	fmt.Println("  $ git unmigrate -de zip")
	fmt.Println("  DRY RUN: git lfs untrack *.zip **/*.zip")
	fmt.Println()
	fmt.Println("  $ git unmigrate -dce mp3")
	fmt.Println("  DRY RUN: git lfs untrack *.mp3 *.MP3 **/*.mp3 **/*.MP3")
	fmt.Println()
	fmt.Println("See also: git-ls-files, git-lfs-track, and git-lfs-untrack")
}

func runGitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
