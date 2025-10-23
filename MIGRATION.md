# Migration from Bash to Go

This document describes the conversion of the Git LFS Scripts project from bash to Go.

## Overview

All bash scripts in the `bin/` directory have been converted to Go commands with the `git-` prefix. The conversion provides better cross-platform compatibility, type safety, and maintainability while preserving all original functionality.

## Architecture

**Chosen Approach:** Independent binaries (Option A)

Each command is a standalone Go binary that can be invoked directly as a Git subcommand:

- `git ls-files` → `cmd/git-ls-files/main.go`
- `git lfs-track` → `cmd/git-lfs-track/main.go`
- etc.

All commands share common code through the `internal/` packages:

- `internal/common/` - Shared utilities (git operations, error handling)
- `internal/lfsfiles/` - Pattern permutation logic (used by ls-files, track, untrack commands)
- `internal/github/` - GitHub API operations (used by delete-github-repo)

## Command Mapping

### Original Bash → New Go Commands

| Bash Script | Go Command | Description |
|------------|------------|-------------|
| `bin/ls-files` | `git-ls-files`, `git-lfs-files`, `git-lfs-track`, `git-lfs-untrack` | Split into 4 separate binaries sharing code |
| `bin/git_lfs_trace` | `git-lfs-trace` | Ruby → Go conversion |
| `bin/nonlfs` | `git-nonlfs` | Bash → Go conversion |
| `bin/unmigrate` | `git-unmigrate` | Bash → Go conversion |
| `bin/new_bare_repo` | `git-new-bare-repo` | Bash → Go conversion |
| `bin/delete_github_repo` | `git-delete-github-repo` | Bash → Go conversion |
| `bin/giftless` | `git-giftless` | Go wrapper for Python/uwsgi |
| `bin/start_lfs_server` | *(not converted)* | Specific to lfs-test-server setup |

## Key Design Decisions

### 1. Multi-name Script Handling

**Question:** How to handle `ls-files` script that responded to multiple names via symlinks?

**Decision:** Create separate binaries (`git-ls-files`, `git-lfs-files`, `git-lfs-track`, `git-lfs-untrack`) that share code via `internal/lfsfiles` package.

**Rationale:**

- Simpler than symlink detection
- Each command has clear purpose
- No runtime overhead checking argv[0]
- Easier to maintain and test

### 2. Flag Parsing

**Change:** Bash used combined flags like `-dce`, Go uses separate flags `-d -c -e`

**Rationale:** Go's standard `flag` package doesn't support GNU-style combined short flags by default. Using separate flags is more idiomatic in Go and clearer.

### 3. Git Subcommand Naming

**Decision:** All commands use `git-` prefix consistently

**Examples:**

- `git-ls-files` (can be invoked as `git ls-files`)
- `git-nonlfs` (can be invoked as `git nonlfs`)
- `git-lfs-track` (can be invoked as `git lfs-track`)

### 4. Giftless Server
**Question:** How to handle Python WSGI application?

**Decision:** Go wrapper that launches Python/uwsgi process

**Rationale:**

- WSGI is Python-specific, no Go equivalent
- Go wrapper provides better process management
- Configuration and validation can be done in Go
- Maintains compatibility with existing Giftless installations

## Installation Methods

### 1. Using Make (Recommended)

```bash
make build
make install                              # Installs to ~/.local/bin
make install INSTALL_DIR=/usr/local/bin   # Custom location
```

### 2. Using Go Install

```bash
go install github.com/mslinn/git_lfs_scripts/cmd/git-ls-files@latest
# repeat for each command
```

### 3. Manual Build

```bash
go build -o build/git-ls-files ./cmd/git-ls-files
# repeat for each command
```

## Testing the Migration

### Build Verification

```bash
make build
ls -lh build/  # Should show all 10 binaries
```

### Functional Testing

```bash
# Test help output
./build/git-ls-files -h
./build/git-nonlfs -h

# Test pattern expansion
./build/git-lfs-track -d -c -e mp3 mp4

# Should output:
# DRY RUN: git lfs track *.mp3 *.MP3 **/*.mp3 **/*.MP3
# DRY RUN: git lfs track *.mp4 *.MP4 **/*.mp4 **/*.MP4
```

## Differences from Original Bash Scripts

### Functional Changes

1. **Flag syntax:** `-dce` → `-d -c -e` (separate flags)
2. **Error handling:** More detailed error messages with proper exit codes
3. **Symlinks:** `ls-files` script is now 4 separate binaries instead of symlink-based dispatch

### Improvements

1. **Cross-platform:** Works on Windows, macOS, Linux without modification
2. **Type safety:** Compile-time checking prevents many runtime errors
3. **Performance:** Go binaries are faster than bash scripts
4. **Maintainability:** Shared code in `internal/` packages reduces duplication
5. **Testing:** Can write Go unit tests for all functionality
6. **Dependencies:** No dependency on bash, find, grep, sed, etc.

## Backward Compatibility

**Note:** Backward compatibility with bash scripts was explicitly NOT required.

If needed in the future:

- Keep bash scripts in `bin/` directory
- Add both directories to PATH
- Users can choose which version to use based on preference

## Future Enhancements

Possible improvements now that we're in Go:

1. **Unit tests:** Add comprehensive test coverage
2. **CI/CD:** Automated builds for multiple platforms
3. **Configuration file:** Add support for `.git-lfs-scripts.yaml`
4. **Parallel operations:** Use goroutines for faster processing
5. **Progress bars:** Add visual feedback for long operations
6. **Native LFS server:** Implement a pure Go LFS server instead of wrapping Python

## Migration Checklist

- [x] Create Go module structure
- [x] Implement `internal/common` package
- [x] Implement `internal/lfsfiles` package for pattern permutation
- [x] Convert `ls-files` → `git-ls-files`, `git-lfs-files`, `git-lfs-track`, `git-lfs-untrack`
- [x] Convert `git_lfs_trace` (Ruby) → `git-lfs-trace` (Go)
- [x] Convert `nonlfs` → `git-nonlfs`
- [x] Convert `unmigrate` → `git-unmigrate`
- [x] Convert `new_bare_repo` → `git-new-bare-repo`
- [x] Implement `internal/github` package
- [x] Convert `delete_github_repo` → `git-delete-github-repo`
- [x] Create `git-giftless` Go wrapper
- [x] Create Makefile for building and installation
- [x] Update README.md with Go installation instructions
- [x] Update .gitignore for Go artifacts
- [x] Build and test all commands
- [x] Document migration process

## Success Criteria

All criteria met:

- ✅ All 10 commands build successfully
- ✅ Help output works for all commands
- ✅ Pattern expansion works correctly (tested with dry-run)
- ✅ All commands follow `git-` naming convention
- ✅ Shared code properly organized in `internal/` packages
- ✅ Makefile provides easy build and install process
- ✅ README updated with comprehensive documentation
- ✅ Standard Go project structure with proper module definition
