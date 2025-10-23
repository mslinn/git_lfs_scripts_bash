package lfsfiles

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/lithammer/dedent"
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
	title := ""
	switch cmdType {
	case LsFiles:
		cmdName = "git-ls-files"
		title = "git-ls-files - Frontend for git ls-files with pattern permutation"
	case LfsLsFiles:
		cmdName = "git-lfs-files"
		title = "git-lfs-files - Frontend for git lfs ls-files with pattern permutation"
	case LfsTrack:
		cmdName = "git-lfs-track"
		title = "git-lfs-track - Frontend for git lfs track with pattern permutation"
	case LfsUntrack:
		cmdName = "git-lfs-untrack"
		title = "git-lfs-untrack - Frontend for git lfs untrack with pattern permutation"
	}

	gitCmd := GetCommandString(cmdType)

	// Build description based on command type
	var helpText string
	if cmdType == LsFiles {
		helpText = dedent.Dedent(fmt.Sprintf(`
			%s

			USAGE:
			  %s [OPTIONS] PATTERN ...

			OPTIONS:
			  -c  Expand pattern to upper and lower case, helpful for media files
			  -d  Dry run (display filename patterns that would be affected)
			  -e  Apply the pattern everywhere (all directories in the Git repository)
			  -h  Show this help message

			DESCRIPTION:
			  This command acts as a frontend to 'git ls-files', permutating wildmatch
			  patterns into more general git ignore/git lfs patterns.

			  Related commands are available for:
			    - git lfs ls-files (git-lfs-files)
			    - git lfs track (git-lfs-track)
			    - git lfs untrack (git-lfs-untrack)

			  Note: These frontend scripts do not support all options of the underlying
			  commands. You can run the underlying commands directly when needed.

			EXAMPLES:
			  # Single pattern dry run
			  %s -d zip
			  # Output: DRY RUN: %s *.zip

			  # Multiple patterns
			  %s -d pdf zip
			  # Output: DRY RUN: %s *.pdf
			  #         DRY RUN: %s *.zip

			  # Case variations
			  %s -dc mp3
			  # Output: DRY RUN: %s *.mp3 *.MP3

			  # Multiple patterns with case variations
			  %s -dc mp3 mp4
			  # Output: DRY RUN: %s *.mp3 *.MP3
			  #         DRY RUN: %s *.mp4 *.MP4

			  # Apply everywhere in repository
			  %s -de zip
			  # Output: DRY RUN: %s *.zip **/*.zip

			  # Combined: everywhere + case variations
			  %s -dce mp3
			  # Output: DRY RUN: %s *.mp3 *.MP3 **/*.mp3 **/*.MP3

			  # Multiple patterns with all options
			  %s -dce mp3 mp4
			  # Output: DRY RUN: %s *.mp3 *.MP3 **/*.mp3 **/*.MP3
			  #         DRY RUN: %s *.mp4 *.MP4 **/*.mp4 **/*.MP4

			SEE ALSO:
			  Related commands: git-lfs-files, git-ls-files, git-lfs-track, git-unmigrate, git-lfs-untrack
			  Documentation: https://mslinn.com/git/5300-git-lfs-patterns-tracking.html
			`, title, cmdName,
			cmdName, gitCmd,
			cmdName, gitCmd, gitCmd,
			cmdName, gitCmd,
			cmdName, gitCmd, gitCmd,
			cmdName, gitCmd,
			cmdName, gitCmd,
			cmdName, gitCmd, gitCmd))
	} else {
		helpText = dedent.Dedent(fmt.Sprintf(`
			%s

			USAGE:
			  %s [OPTIONS] PATTERN ...

			OPTIONS:
			  -c  Expand pattern to upper and lower case, helpful for media files
			  -d  Dry run (display filename patterns that would be affected)
			  -e  Apply the pattern everywhere (all directories in the Git repository)
			  -h  Show this help message

			DESCRIPTION:
			  This command permutates wildmatch patterns for use with the underlying
			  Git or Git LFS command.

			EXAMPLES:
			  # Single pattern dry run
			  %s -d zip
			  # Output: DRY RUN: %s *.zip

			  # Multiple patterns
			  %s -d pdf zip
			  # Output: DRY RUN: %s *.pdf
			  #         DRY RUN: %s *.zip

			  # Case variations
			  %s -dc mp3
			  # Output: DRY RUN: %s *.mp3 *.MP3

			  # Multiple patterns with case variations
			  %s -dc mp3 mp4
			  # Output: DRY RUN: %s *.mp3 *.MP3
			  #         DRY RUN: %s *.mp4 *.MP4

			  # Apply everywhere in repository
			  %s -de zip
			  # Output: DRY RUN: %s *.zip **/*.zip

			  # Combined: everywhere + case variations
			  %s -dce mp3
			  # Output: DRY RUN: %s *.mp3 *.MP3 **/*.mp3 **/*.MP3

			  # Multiple patterns with all options
			  %s -dce mp3 mp4
			  # Output: DRY RUN: %s *.mp3 *.MP3 **/*.mp3 **/*.MP3
			  #         DRY RUN: %s *.mp4 *.MP4 **/*.mp4 **/*.MP4

			SEE ALSO:
			  Related commands: git-lfs-files, git-ls-files, git-lfs-track, git-unmigrate, git-lfs-untrack
			  Documentation: https://mslinn.com/git/5300-git-lfs-patterns-tracking.html
			`, title, cmdName,
			cmdName, gitCmd,
			cmdName, gitCmd, gitCmd,
			cmdName, gitCmd,
			cmdName, gitCmd, gitCmd,
			cmdName, gitCmd,
			cmdName, gitCmd,
			cmdName, gitCmd, gitCmd))
	}

	fmt.Print(helpText)
}
