.PHONY: all build install uninstall clean test help build-release-tool

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install

# Binary names
BINARIES := \
	git-ls-files \
	git-lfs-files \
	git-lfs-track \
	git-lfs-untrack \
	git-lfs-trace \
	git-nonlfs \
	git-unmigrate \
	git-new-bare-repo \
	git-delete-github-repo \
	git-giftless

# Build directory
BUILD_DIR := build

# Install directory using Go's standard mechanism
# GOBIN if set, otherwise GOPATH/bin, otherwise ~/go/bin as fallback
GOBIN ?= $(shell go env GOBIN)
ifeq ($(GOBIN),)
	GOPATH ?= $(shell go env GOPATH)
	ifeq ($(GOPATH),)
		INSTALL_DIR := $(HOME)/go/bin
	else
		INSTALL_DIR := $(GOPATH)/bin
	endif
else
	INSTALL_DIR := $(GOBIN)
endif

# Version info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/mslinn/git_lfs_scripts/internal/common.Version=$(VERSION)"

all: build

help:
	@echo "Git LFS Scripts - Makefile targets:"
	@echo ""
	@echo "  make build         Build all binaries to $(BUILD_DIR)/"
	@echo "  make install       Install binaries to GOBIN or GOPATH/bin"
	@echo "  make uninstall     Remove installed binaries"
	@echo "  make clean         Remove built binaries"
	@echo "  make test          Run tests"
	@echo "  make help          Show this help message"
	@echo ""
	@echo "Environment variables:"
	@echo "  GOBIN              Go binary installation directory (overrides GOPATH/bin)"
	@echo "  GOPATH             Go workspace path (default: ~/go)"
	@echo "  VERSION            Version string (default: git describe or 'dev')"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make install"
	@echo "  GOBIN=/usr/local/bin make install"

build:
	@mkdir -p $(BUILD_DIR)
	@echo "Building all binaries..."
	@for bin in $(BINARIES); do \
		echo "  Building $$bin..."; \
		$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$$bin ./cmd/$$bin || exit 1; \
	done
	@echo "Build complete! Binaries are in $(BUILD_DIR)/"

install: build
	@echo "Installing binaries to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@for bin in $(BINARIES); do \
		echo "  Installing $$bin..."; \
		/usr/bin/install -m 755 $(BUILD_DIR)/$$bin $(INSTALL_DIR)/$$bin || exit 1; \
	done
	@echo "Installation complete!"
	@echo ""
	@echo "Make sure $(INSTALL_DIR) is in your PATH."
	@echo ""
	@echo "Installed Git LFS helper commands:"
	@echo "  git ls-files           - List files with pattern expansion"
	@echo "  git lfs-files          - List Git LFS tracked files with pattern expansion"
	@echo "  git lfs-track          - Track patterns in Git LFS with expansion"
	@echo "  git lfs-untrack        - Untrack patterns from Git LFS with expansion"
	@echo "  git lfs-trace          - Git LFS transfer adapter for debugging"
	@echo "  git nonlfs             - List files NOT in Git LFS"
	@echo "  git unmigrate          - Reverse 'git lfs migrate import'"
	@echo "  git new-bare-repo      - Create new bare Git repositories"
	@echo "  git delete-github-repo - Delete GitHub repositories (requires gh CLI)"
	@echo "  git giftless           - Go wrapper for Python Giftless LFS server"

uninstall: ## Remove installed binaries
	@echo "Uninstalling binaries from $(INSTALL_DIR)..."
	@for bin in $(BINARIES); do \
		if [ -f $(INSTALL_DIR)/$$bin ]; then \
			echo "  Removing $$bin..."; \
			rm -f $(INSTALL_DIR)/$$bin; \
		fi \
	done
	@echo "Uninstall complete!"

clean:
	@echo "Cleaning build artifacts..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f release
	@echo "Clean complete!"

test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

# Target to just download dependencies
deps:
	@echo "Downloading dependencies..."
	@$(GOCMD) mod download
	@echo "Dependencies downloaded!"

# Target to tidy up go.mod and go.sum
tidy:
	@echo "Tidying go.mod..."
	@$(GOCMD) mod tidy
	@echo "Tidy complete!"

# Build for specific OS/architecture
build-linux:
	@GOOS=linux GOARCH=amd64 $(MAKE) build

build-darwin:
	@GOOS=darwin GOARCH=amd64 $(MAKE) build

build-windows:
	@GOOS=windows GOARCH=amd64 $(MAKE) build

build-release-tool: ## Build the release tool (developers only)
	@echo "Building release tool (for developers)..."
	$(GOBUILD) -o release ./cmd/release
	@echo "Build complete: ./release"
	@echo "Note: This is a development tool and is not installed with 'go install'"
