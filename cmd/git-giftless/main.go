package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/lithammer/dedent"
	"github.com/mslinn/git_lfs_scripts/internal/common"
)

const (
	defaultVenvPath = "/opt/giftless/.venv/bin/activate"
	defaultHost     = "0.0.0.0"
	defaultPort     = "9876"
)

func main() {
	var (
		venvPath string
		host     string
		port     string
		threads  int
		workers  int
		showHelp bool
	)

	flag.StringVar(&venvPath, "venv", defaultVenvPath, "Path to Python virtual environment activation script")
	flag.StringVar(&host, "host", defaultHost, "Host address to bind to")
	flag.StringVar(&port, "port", defaultPort, "Port to listen on")
	flag.IntVar(&threads, "threads", 2, "Number of threads per worker")
	flag.IntVar(&workers, "workers", 2, "Number of worker processes")
	flag.BoolVar(&showHelp, "h", false, "Show help")
	flag.Parse()

	if showHelp {
		printHelp()
		os.Exit(0)
	}

	// Check all prerequisites before starting
	checkPrerequisites()

	fmt.Printf("Starting Giftless LFS server on %s:%s\n", host, port)
	fmt.Printf("Workers: %d, Threads: %d\n", workers, threads)

	// Build uwsgi command
	cmd := exec.Command("uwsgi",
		"--master",
		fmt.Sprintf("--threads=%d", threads),
		fmt.Sprintf("--processes=%d", workers),
		"--manage-script-name",
		"--module=giftless.wsgi_entrypoint",
		"--callable=app",
		fmt.Sprintf("--http=%s:%s", host, port),
	)

	// If venv path exists, we need to activate it first
	// For simplicity, we'll use bash to source the venv and run uwsgi
	if _, err := os.Stat(venvPath); err == nil {
		bashCmd := fmt.Sprintf("source %s && uwsgi --master --threads=%d --processes=%d --manage-script-name --module=giftless.wsgi_entrypoint --callable=app --http=%s:%s",
			venvPath, threads, workers, host, port)

		cmd = exec.Command("bash", "-c", bashCmd)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Handle signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- cmd.Run()
	}()

	// Wait for either completion or signal
	select {
	case err := <-errChan:
		if err != nil {
			common.PrintError("Server exited with error: %v", err)
		}
		fmt.Println("Server stopped")
	case sig := <-sigChan:
		fmt.Printf("\nReceived signal %v, shutting down...\n", sig)
		if cmd.Process != nil {
			cmd.Process.Signal(sig)
		}
		// Wait for process to exit
		<-errChan
		fmt.Println("Server stopped")
	}
}

func printHelp() {
	fmt.Print(dedent.Dedent(`
		git-giftless - Start a Giftless Git LFS server

		USAGE:
		  git giftless [OPTIONS]

		OPTIONS:
		  --venv PATH      Path to Python virtual environment (default: /opt/giftless/.venv/bin/activate)
		  --host ADDRESS   Host address to bind to (default: 0.0.0.0)
		  --port PORT      Port to listen on (default: 9876)
		  --threads N      Number of threads per worker (default: 2)
		  --workers N      Number of worker processes (default: 2)
		  -h, --help       Show this help message

		DESCRIPTION:
		  This command starts a Giftless Git LFS server using uwsgi as a WSGI server.
		  All prerequisites are verified before starting the server.

		REQUIREMENTS:
		  - Python 3 (python3 command must be available)
		  - Giftless direct dependencies:
		    azure-storage-blob, boto3, cachetools, cryptography, figcan,
		    flask, flask-classful, flask-marshmallow, google-cloud-storage,
		    importlib-metadata, pyjwt, python-dateutil, python-dotenv,
		    pyyaml, typing-extensions, webargs, werkzeug
		  - giftless Python package (pip install giftless)
		  - uwsgi Python package (pip install uwsgi)

		EXAMPLES:
		  # Start with defaults
		  git giftless

		  # Custom port and workers
		  git giftless --port 8080 --workers 4

		  # Use specific virtual environment
		  git giftless --venv /path/to/venv/bin/activate
	`))
}

func checkPrerequisites() {
	var missing []string
	var missingPackages []string

	// Check Python 3
	if err := checkCommand("python3", "--version"); err != nil {
		missing = append(missing, "Python 3 (install from: https://www.python.org/)")
	}

	// Check giftless direct dependencies
	deps := []struct {
		module string
		pkg    string
	}{
		{"azure.storage.blob", "azure-storage-blob"},
		{"boto3", "boto3"},
		{"cachetools", "cachetools"},
		{"cryptography", "cryptography"},
		{"figcan", "figcan"},
		{"flask", "flask"},
		{"flask_classful", "flask-classful"},
		{"flask_marshmallow", "flask-marshmallow"},
		{"google.cloud.storage", "google-cloud-storage"},
		{"importlib_metadata", "importlib-metadata"},
		{"jwt", "pyjwt"},
		{"dateutil", "python-dateutil"},
		{"dotenv", "python-dotenv"},
		{"yaml", "pyyaml"},
		{"typing_extensions", "typing-extensions"},
		{"webargs", "webargs"},
		{"werkzeug", "werkzeug"},
	}

	for _, dep := range deps {
		if err := checkCommand("python3", "-c", "import "+dep.module); err != nil {
			missing = append(missing, dep.pkg)
			missingPackages = append(missingPackages, dep.pkg)
		}
	}

	// Check giftless
	if err := checkCommand("python3", "-c", "import giftless"); err != nil {
		missing = append(missing, "giftless")
		missingPackages = append(missingPackages, "giftless")
	}

	// Check uwsgi
	if err := checkCommand("uwsgi", "--version"); err != nil {
		missing = append(missing, "uwsgi")
		missingPackages = append(missingPackages, "uwsgi")
	}

	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "Error: Missing required dependencies:\n")
		for _, dep := range missing {
			fmt.Fprintf(os.Stderr, "  ✗ %s\n", dep)
		}
		fmt.Fprintf(os.Stderr, "\nTo install all missing dependencies, run:\n")
		fmt.Fprintf(os.Stderr, "  pip install %s\n", strings.Join(missingPackages, " "))
		os.Exit(1)
	}

	fmt.Println("✓ All prerequisites verified")
}

func checkCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	// Suppress output, we only care about exit code
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command '%s' not found or failed", name)
	}
	return nil
}
