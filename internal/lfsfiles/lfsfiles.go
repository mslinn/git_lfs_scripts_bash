package lfsfiles

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// CommandType represents the type of git command to execute
type CommandType int

const (
	LsFiles CommandType = iota
	LfsLsFiles
	LfsTrack
	LfsUntrack
)

// Options holds the command-line options
type Options struct {
	BothCases  bool   // -c: Expand pattern to upper and lower case
	DryRun     bool   // -d: Dry run
	Everywhere bool   // -e: Apply pattern everywhere (all directories)
	Command    string // The git command to execute
}

// ExpandPattern expands a file extension pattern based on options
func ExpandPattern(pattern string, opts Options) []string {
	var patterns []string

	lc := strings.ToLower(pattern)
	uc := strings.ToUpper(pattern)

	if opts.Everywhere {
		if opts.BothCases {
			patterns = []string{
				"*." + lc,
				"*." + uc,
				"**/*." + lc,
				"**/*." + uc,
			}
		} else {
			patterns = []string{
				"*." + pattern,
				"**/*." + pattern,
			}
		}
	} else {
		if opts.BothCases {
			patterns = []string{
				"*." + lc,
				"*." + uc,
			}
		} else {
			patterns = []string{
				"*." + pattern,
			}
		}
	}

	return patterns
}

// Execute runs the git command with expanded patterns
func Execute(patterns []string, opts Options) error {
	if opts.DryRun {
		for _, pattern := range patterns {
			expanded := ExpandPattern(pattern, opts)
			fmt.Printf("DRY RUN: %s %s\n", opts.Command, strings.Join(expanded, " "))
		}
		return nil
	}

	// If no patterns provided and it's a ls-files command, just run the command
	if len(patterns) == 0 && (opts.Command == "git ls-files" || opts.Command == "git lfs ls-files") {
		return executeCommand(opts.Command, []string{})
	}

	// Execute command for each pattern
	for _, pattern := range patterns {
		expanded := ExpandPattern(pattern, opts)
		if err := executeCommand(opts.Command, expanded); err != nil {
			return err
		}
	}

	return nil
}

// executeCommand runs a git command with the given arguments
func executeCommand(cmdStr string, args []string) error {
	parts := strings.Fields(cmdStr)
	allArgs := append(parts[1:], args...)

	cmd := exec.Command(parts[0], allArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// GetCommandString returns the git command string for the given command type
func GetCommandString(cmdType CommandType) string {
	switch cmdType {
	case LsFiles:
		return "git ls-files"
	case LfsLsFiles:
		return "git lfs ls-files"
	case LfsTrack:
		return "git lfs track"
	case LfsUntrack:
		return "git lfs untrack"
	default:
		return ""
	}
}

// PrintHelp prints help message for the command
func PrintHelp(cmdType CommandType) {
	cmdName := ""
	switch cmdType {
	case LsFiles:
		cmdName = "git-ls-files"
	case LfsLsFiles:
		cmdName = "git-lfs-files"
	case LfsTrack:
		cmdName = "git-lfs-track"
	case LfsUntrack:
		cmdName = "git-lfs-untrack"
	}

	gitCmd := GetCommandString(cmdType)

	description := ""
	switch cmdType {
	case LsFiles:
		description = `git-ls-files - Invoke git ls-files with permutated wildmatch patterns.

In addition to acting as a front-end to 'git ls-files', there are related
commands for 'git lfs ls-files', 'git lfs track' and 'git lfs untrack'.

The front-end scripts permutate wildmatch patterns into a more general
git ignore / git lfs pattern. They do not support all options supported by
the commands that they invoke. You can, of course, run those commands when required.`
	case LfsLsFiles:
		description = "git-lfs-files - Invoke git lfs ls-files with permutated wildmatch patterns."
	case LfsTrack:
		description = "git-lfs-track - Invoke git lfs track with permutated wildmatch patterns."
	case LfsUntrack:
		description = "git-lfs-untrack - Invoke git lfs untrack with permutated wildmatch patterns."
	}

	fmt.Printf("%s\n\n", description)
	fmt.Printf("Syntax:\n")
	fmt.Printf("  %s [OPTIONS] PATTERN ...\n\n", cmdName)
	fmt.Printf("OPTIONS:\n")
	fmt.Printf("  -c  Expand pattern to upper and lower case, helpful for media files\n")
	fmt.Printf("  -d  Dry run (display filename patterns that would be affected)\n")
	fmt.Printf("  -e  Apply the pattern everywhere (all directories in the Git repository)\n")
	fmt.Printf("  -h  Show this help message\n\n")
	fmt.Printf("Examples:\n")
	fmt.Printf("  $ %s -d zip\n", cmdName)
	fmt.Printf("  DRY RUN: %s *.zip\n\n", gitCmd)
	fmt.Printf("  $ %s -d pdf zip\n", cmdName)
	fmt.Printf("  DRY RUN: %s *.pdf\n", gitCmd)
	fmt.Printf("  DRY RUN: %s *.zip\n\n", gitCmd)
	fmt.Printf("  $ %s -dc mp3\n", cmdName)
	fmt.Printf("  DRY RUN: %s *.mp3 *.MP3\n\n", gitCmd)
	fmt.Printf("  $ %s -dc mp3 mp4\n", cmdName)
	fmt.Printf("  DRY RUN: %s *.mp3 *.MP3\n", gitCmd)
	fmt.Printf("  DRY RUN: %s *.mp4 *.MP4\n\n", gitCmd)
	fmt.Printf("  $ %s -de zip\n", cmdName)
	fmt.Printf("  DRY RUN: %s *.zip **/*.zip\n\n", gitCmd)
	fmt.Printf("  $ %s -dce mp3\n", cmdName)
	fmt.Printf("  DRY RUN: %s *.mp3 *.MP3 **/*.mp3 **/*.MP3\n\n", gitCmd)
	fmt.Printf("  $ %s -dce mp3 mp4\n", cmdName)
	fmt.Printf("  DRY RUN: %s *.mp3 *.MP3 **/*.mp3 **/*.MP3\n", gitCmd)
	fmt.Printf("  DRY RUN: %s *.mp4 *.MP4 **/*.mp4 **/*.MP4\n\n", gitCmd)
	fmt.Printf("See also:\n")
	fmt.Printf("  Related commands: git-lfs-files, git-ls-files, git-lfs-track, git-unmigrate, and git-lfs-untrack.\n")
	fmt.Printf("  https://mslinn.com/git/5300-git-lfs-patterns-tracking.html\n")
}
