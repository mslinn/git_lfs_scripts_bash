package main

import (
	"bufio"
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

	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	// Check if we're in a git repository
	if err := common.CheckGitRepo(); err != nil {
		common.PrintError("%v", err)
	}

	// Get all files in the repository (excluding .git directory)
	allFiles, err := getAllFiles()
	if err != nil {
		common.PrintError("Failed to get all files: %v", err)
	}

	// Get LFS tracked patterns from .gitattributes
	lfsPatterns, err := getLFSPatterns()
	if err != nil {
		common.PrintError("Failed to get LFS patterns: %v", err)
	}

	// Find files matching LFS patterns
	lfsFiles := make(map[string]bool)
	for _, pattern := range lfsPatterns {
		matches, _ := findMatchingFiles(pattern)
		for _, match := range matches {
			lfsFiles[match] = true
		}
	}

	// Print files that are NOT in LFS
	for _, file := range allFiles {
		if !lfsFiles[file] {
			fmt.Println(file)
		}
	}
}

func printHelp() {
	fmt.Print(dedent.Dedent(`
		git-nonlfs - List files that are not managed by Git LFS

		USAGE:
		  git nonlfs [OPTIONS]

		OPTIONS:
		  -h  Show this help message

		DESCRIPTION:
		  This command lists all files in the repository that are not tracked by Git LFS.
		  It reads .gitattributes to determine which patterns are tracked by LFS, then
		  lists all files that don't match those patterns.

		  Requires:
		    - Git repository
		    - find command (standard on Unix/Linux/macOS)

		EXAMPLES:
		  # List all non-LFS files
		  git nonlfs

		  # Count non-LFS files
		  git nonlfs | wc -l

		  # Find large non-LFS files
		  git nonlfs | xargs du -h | sort -hr | head -10
	`))
}

func getAllFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .git directory
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// Only include files, not directories
		if !info.IsDir() {
			// Remove leading "./"
			cleanPath := strings.TrimPrefix(path, "./")
			files = append(files, cleanPath)
		}

		return nil
	})

	return files, err
}

func getLFSPatterns() ([]string, error) {
	file, err := os.Open(".gitattributes")
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil // No .gitattributes file
		}
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse lines like "*.pdf filter=lfs diff=lfs merge=lfs -text"
		fields := strings.Fields(line)
		if len(fields) > 0 && strings.Contains(line, "filter=lfs") {
			patterns = append(patterns, fields[0])
		}
	}

	return patterns, scanner.Err()
}

func findMatchingFiles(pattern string) ([]string, error) {
	// Use find command to locate files matching the pattern
	cmd := exec.Command("find", ".", "-name", pattern, "-type", "f")
	output, err := cmd.Output()
	if err != nil {
		return []string{}, nil
	}

	var files []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		path := strings.TrimPrefix(scanner.Text(), "./")
		files = append(files, path)
	}

	return files, nil
}
