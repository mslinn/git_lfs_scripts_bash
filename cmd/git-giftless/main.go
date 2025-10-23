package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

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

	// Check if uwsgi is available
	if err := checkCommand("uwsgi"); err != nil {
		common.PrintError("uwsgi is not installed. Install with: pip install uwsgi")
	}

	// Check if giftless is available
	if err := checkCommand("python3", "-c", "import giftless"); err != nil {
		common.PrintError("giftless is not installed. Install with: pip install giftless")
	}

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
	fmt.Println("git-giftless - Start a Giftless Git LFS server")
	fmt.Println()
	fmt.Println("Usage: git giftless [OPTIONS]")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  --venv PATH      Path to Python virtual environment (default: /opt/giftless/.venv/bin/activate)")
	fmt.Println("  --host ADDRESS   Host address to bind to (default: 0.0.0.0)")
	fmt.Println("  --port PORT      Port to listen on (default: 9876)")
	fmt.Println("  --threads N      Number of threads per worker (default: 2)")
	fmt.Println("  --workers N      Number of worker processes (default: 2)")
	fmt.Println("  -h, --help       Show this help message")
	fmt.Println()
	fmt.Println("This command starts a Giftless Git LFS server using uwsgi.")
	fmt.Println()
	fmt.Println("Requirements:")
	fmt.Println("  - Python 3 with giftless installed (pip install giftless)")
	fmt.Println("  - uwsgi installed (pip install uwsgi)")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  git giftless --port 8080 --workers 4")
}

func checkCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command '%s' not found or failed", name)
	}
	return nil
}
