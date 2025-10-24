package lfsfiles

import (
	"reflect"
	"testing"
)

// TestExpandPattern tests the wildmatch pattern expansion logic
func TestExpandPattern(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		opts     Options
		expected []string
	}{
		{
			name:    "basic pattern - current directory only",
			pattern: "zip",
			opts: Options{
				BothCases:  false,
				Everywhere: false,
			},
			expected: []string{"*.zip"},
		},
		{
			name:    "case variations - current directory",
			pattern: "mp3",
			opts: Options{
				BothCases:  true,
				Everywhere: false,
			},
			expected: []string{"*.mp3", "*.MP3"},
		},
		{
			name:    "everywhere - single case",
			pattern: "pdf",
			opts: Options{
				BothCases:  false,
				Everywhere: true,
			},
			expected: []string{"*.pdf", "**/*.pdf"},
		},
		{
			name:    "everywhere with case variations",
			pattern: "mp4",
			opts: Options{
				BothCases:  true,
				Everywhere: true,
			},
			expected: []string{"*.mp4", "*.MP4", "**/*.mp4", "**/*.MP4"},
		},
		{
			name:    "uppercase input - case variations",
			pattern: "JPG",
			opts: Options{
				BothCases:  true,
				Everywhere: false,
			},
			expected: []string{"*.jpg", "*.JPG"},
		},
		{
			name:    "mixed case input - case variations everywhere",
			pattern: "MoV",
			opts: Options{
				BothCases:  true,
				Everywhere: true,
			},
			expected: []string{"*.mov", "*.MOV", "**/*.mov", "**/*.MOV"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandPattern(tt.pattern, tt.opts)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExpandPattern(%q, %+v) = %v, want %v",
					tt.pattern, tt.opts, result, tt.expected)
			}
		})
	}
}

// TestExpandPatternOrder tests that patterns are in the correct order
// Order matters for git commands to work properly
func TestExpandPatternOrder(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		opts     Options
		expected []string
	}{
		{
			name:    "everywhere order - current dir before subdirs",
			pattern: "txt",
			opts: Options{
				BothCases:  false,
				Everywhere: true,
			},
			expected: []string{"*.txt", "**/*.txt"},
		},
		{
			name:    "case order - lowercase before uppercase",
			pattern: "log",
			opts: Options{
				BothCases:  true,
				Everywhere: false,
			},
			expected: []string{"*.log", "*.LOG"},
		},
		{
			name:    "combined order - current lower, current upper, subdir lower, subdir upper",
			pattern: "dat",
			opts: Options{
				BothCases:  true,
				Everywhere: true,
			},
			expected: []string{"*.dat", "*.DAT", "**/*.dat", "**/*.DAT"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandPattern(tt.pattern, tt.opts)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Pattern order mismatch:\ngot:  %v\nwant: %v",
					result, tt.expected)
			}
		})
	}
}

// TestGetCommandString tests command type to string conversion
func TestGetCommandString(t *testing.T) {
	tests := []struct {
		name     string
		cmdType  CommandType
		expected string
	}{
		{
			name:     "ls-files command",
			cmdType:  LsFiles,
			expected: "git ls-files",
		},
		{
			name:     "lfs ls-files command",
			cmdType:  LfsLsFiles,
			expected: "git lfs ls-files",
		},
		{
			name:     "lfs track command",
			cmdType:  LfsTrack,
			expected: "git lfs track",
		},
		{
			name:     "lfs untrack command",
			cmdType:  LfsUntrack,
			expected: "git lfs untrack",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCommandString(tt.cmdType)
			if result != tt.expected {
				t.Errorf("GetCommandString(%v) = %q, want %q",
					tt.cmdType, result, tt.expected)
			}
		})
	}
}

// TestMediaFilePatterns tests common media file extensions
// These require case variations due to FAT32 filesystem conventions
func TestMediaFilePatterns(t *testing.T) {
	mediaExtensions := []string{"mp3", "mp4", "mov", "avi", "jpg", "png", "gif"}

	for _, ext := range mediaExtensions {
		t.Run("media_"+ext, func(t *testing.T) {
			opts := Options{
				BothCases:  true,
				Everywhere: true,
			}
			result := ExpandPattern(ext, opts)

			// Should have 4 patterns: current lower, current upper, subdir lower, subdir upper
			if len(result) != 4 {
				t.Errorf("Media file %q should expand to 4 patterns, got %d: %v",
					ext, len(result), result)
			}

			// Verify lowercase and uppercase variants exist
			hasLower := false
			hasUpper := false
			for _, pattern := range result {
				if pattern == "*."+ext || pattern == "**/*."+ext {
					hasLower = true
				}
			}
			// Simple uppercase check - patterns should differ at positions 1 and 3
			if len(result) >= 4 && result[1] != result[0] && result[3] != result[2] {
				hasUpper = true
			}

			if !hasLower {
				t.Errorf("Media file %q missing lowercase patterns in %v", ext, result)
			}
			if !hasUpper {
				t.Errorf("Media file %q missing uppercase patterns in %v", ext, result)
			}
		})
	}
}

// TestWildmatchPatternExamples tests examples from the documentation
func TestWildmatchPatternExamples(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		opts        Options
		description string
		shouldMatch []string // What files this pattern should conceptually match
	}{
		{
			name:    "doc example: current directory only",
			pattern: "zip",
			opts: Options{
				BothCases:  false,
				Everywhere: false,
			},
			description: "Matches *.zip in current directory only",
			shouldMatch: []string{"file.zip", "archive.zip"},
		},
		{
			name:    "doc example: all subdirectories",
			pattern: "zip",
			opts: Options{
				BothCases:  false,
				Everywhere: true,
			},
			description: "Matches *.zip everywhere in repository",
			shouldMatch: []string{"file.zip", "docs/file.zip", "src/data/archive.zip"},
		},
		{
			name:    "doc example: FAT32 media files",
			pattern: "mp3",
			opts: Options{
				BothCases:  true,
				Everywhere: false,
			},
			description: "Handles FAT32 uppercase convention for media files",
			shouldMatch: []string{"song.mp3", "SONG.MP3"},
		},
		{
			name:    "doc example: full combination",
			pattern: "mp4",
			opts: Options{
				BothCases:  true,
				Everywhere: true,
			},
			description: "Case variations everywhere for media files",
			shouldMatch: []string{"video.mp4", "VIDEO.MP4", "clips/video.mp4", "clips/VIDEO.MP4"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandPattern(tt.pattern, tt.opts)
			t.Logf("%s: %v", tt.description, result)
			t.Logf("Should conceptually match: %v", tt.shouldMatch)

			// Verify we got the right number of patterns
			expectedCount := 1
			if tt.opts.Everywhere {
				expectedCount *= 2
			}
			if tt.opts.BothCases {
				expectedCount *= 2
			}

			if len(result) != expectedCount {
				t.Errorf("Expected %d patterns, got %d: %v",
					expectedCount, len(result), result)
			}
		})
	}
}
