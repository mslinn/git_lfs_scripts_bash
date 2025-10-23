.PHONY: all build install clean test help

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

# Install directory (default to user's local bin, can be overridden)
INSTALL_DIR ?= $(HOME)/.local/bin

# Version info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/mslinn/git_lfs_scripts/internal/common.Version=$(VERSION)"

all: build

help:
	@echo "Git LFS Scripts - Makefile targets:"
	@echo ""
	@echo "  make build         Build all binaries to $(BUILD_DIR)/"
	@echo "  make install       Install binaries to $(INSTALL_DIR)"
	@echo "  make clean         Remove built binaries"
	@echo "  make test          Run tests"
	@echo "  make help          Show this help message"
	@echo ""
	@echo "Environment variables:"
	@echo "  INSTALL_DIR        Installation directory (default: ~/.local/bin)"
	@echo "  VERSION            Version string (default: git describe or 'dev')"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make install"
	@echo "  make install INSTALL_DIR=/usr/local/bin"

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
		install -m 755 $(BUILD_DIR)/$$bin $(INSTALL_DIR)/$$bin || exit 1; \
	done
	@echo "Installation complete!"
	@echo ""
	@echo "Make sure $(INSTALL_DIR) is in your PATH."
	@echo "You can now use commands like:"
	@echo "  git ls-files"
	@echo "  git lfs-track"
	@echo "  git nonlfs"
	@echo "  etc."

clean:
	@echo "Cleaning build artifacts..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
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
