package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mslinn/git_lfs_scripts/internal/lfsfiles"
)

func main() {
	var opts lfsfiles.Options
	var showHelp bool

	flag.BoolVar(&opts.BothCases, "c", false, "Expand pattern to upper and lower case")
	flag.BoolVar(&opts.DryRun, "d", false, "Dry run")
	flag.BoolVar(&opts.Everywhere, "e", false, "Apply pattern everywhere")
	flag.BoolVar(&showHelp, "h", false, "Show help")
	flag.Parse()

	if showHelp {
		lfsfiles.PrintHelp(lfsfiles.LsFiles)
		os.Exit(0)
	}

	opts.Command = lfsfiles.GetCommandString(lfsfiles.LsFiles)
	patterns := flag.Args()

	// For ls-files, if no patterns provided, just run the command
	// For track/untrack, patterns are required
	if err := lfsfiles.Execute(patterns, opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
