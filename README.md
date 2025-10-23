# Git LFS Scripts

Generally useful Git subcommands for working with Git LFS.

These tools were written along with the miniseries of articles about
[Git LFS on `mslinn.com`](https://www.mslinn.com/git/5100-git-lfs.html).


## Commands

All commands can be invoked as Git subcommands (e.g., `git ls-files`, `git nonlfs`):

* `git-delete-github-repo` - Deletes the given GitHub repo without prompting (requires `gh` CLI)
* `git-giftless`           - Run Giftless Git LFS server (requires Python with giftless and uwsgi)
* `git-lfs-trace`          - Git LFS transfer adapter that reports activity between Git client and LFS server
* `git-ls-files`           - Frontend for `git ls-files` with pattern permutation
* `git-lfs-files`          - Frontend for `git lfs ls-files` with pattern permutation
* `git-lfs-track`          - Frontend for `git lfs track` with pattern permutation
* `git-lfs-untrack`        - Frontend for `git lfs untrack` with pattern permutation
* `git-new-bare-repo`      - Creates a bare Git repository
* `git-nonlfs`             - Lists files that are not in Git LFS
* `git-unmigrate`          - Reverses `git lfs migrate import` for given wildmatch patterns


## Installation

### Prerequisites

* Go 1.18 or later
* Git
* For `git-giftless`: Python 3 with `giftless` and `uwsgi` installed
* For `git-delete-github-repo`: GitHub CLI (`gh`)

### Build and Install

```shell
# Clone the repository
git clone https://github.com/mslinn/git_lfs_scripts.git
cd git_lfs_scripts

# Build all binaries
make build

# Install to ~/.local/bin (default)
make install

# Or install to a custom location
make install INSTALL_DIR=/usr/local/bin
```

Make sure the installation directory is in your `PATH`. For `~/.local/bin`:

```shell
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### Using Go Install Directly

Alternatively, you can install individual commands directly:

```shell
go install github.com/mslinn/git_lfs_scripts/cmd/git-ls-files@latest
go install github.com/mslinn/git_lfs_scripts/cmd/git-nonlfs@latest
# etc.
```


## Usage Examples

### Pattern Permutation Commands

These commands expand file extension patterns for convenience:

```shell
# Track all PDF files (current directory only)
git lfs-track pdf

# Track all MP3 files in upper and lower case (everywhere)
git lfs-track -ce mp3

# Dry run to see what would be tracked
git lfs-track -dce mp3 mp4

# Flags can be combined or separate
git lfs-track -d -c -e mp3    # Same as -dce

# Long flag names are also supported
git lfs-track --dryrun --bothcases --everywhere mp3

# List all files not tracked by LFS
git nonlfs

# Unmigrate files from LFS back to Git
git unmigrate -ce mp3
```

#### Common Flags

Commands that support pattern permutation (`git-ls-files`, `git-lfs-files`, `git-lfs-track`, `git-lfs-untrack`) support:

* `-c`, `--bothcases` - Expand pattern to upper and lower case (useful for media files)
* `-d`, `--dryrun`     - Show what would be done without executing
* `-e`, `--everywhere` - Apply pattern recursively in all directories
* `-h`, `--help`       - Show help message

Flags can be combined (e.g., `-dce`) or used separately (e.g., `-d -c -e`).

### Server and Repository Commands

```shell
# Start Giftless LFS server
git giftless --port 8080 --workers 4

# Create a new bare repository
git new-bare-repo /path/to/repo.git

# Delete a GitHub repository
git delete-github-repo my-test-repo
```

### LFS Trace Adapter

To use the LFS trace adapter, configure it in your Git LFS config:

```shell
git config lfs.customtransfer.trace.path `which git-lfs-trace`
```

The above adds something similar to the following to the current Git repository configuration
(`.git/config`)

```text
[lfs "customtransfer.trace"]
  path = /home/mslinn/go/bin/git-lfs-trace
```


## Development

### Building

```shell
make build          # Build all binaries
make test           # Run tests
make clean          # Clean build artifacts
make tidy           # Tidy go.mod
```

### Project Structure

```text
.
├── cmd/                    # Command implementations
│   ├── git-ls-files/
│   ├── git-lfs-files/
│   ├── git-lfs-track/
│   ├── git-lfs-untrack/
│   ├── git-lfs-trace/
│   ├── git-nonlfs/
│   ├── git-unmigrate/
│   ├── git-new-bare-repo/
│   ├── git-delete-github-repo/
│   └── git-giftless/
├── internal/               # Shared internal packages
│   ├── common/            # Common utilities
│   ├── lfsfiles/          # Pattern permutation logic
│   └── github/            # GitHub operations
├── Makefile               # Build automation
└── README.md              # This file
```
