package main

import (
	"fmt"
	"os"

	"github.com/mslinn/git_lfs_scripts/internal/lfsfiles"
	"github.com/spf13/pflag"
)

func main() {
	var opts lfsfiles.Options
	var showHelp bool

	pflag.BoolVarP(&opts.BothCases, "bothcases", "c", false, "Expand pattern to upper and lower case")
	pflag.BoolVarP(&opts.DryRun, "dryrun", "d", false, "Dry run")
	pflag.BoolVarP(&opts.Everywhere, "everywhere", "e", false, "Apply pattern everywhere")
	pflag.BoolVarP(&showHelp, "help", "h", false, "Show help")
	pflag.Parse()

	if showHelp {
		lfsfiles.PrintHelp(lfsfiles.LsFiles)
		os.Exit(0)
	}

	opts.Command = lfsfiles.GetCommandString(lfsfiles.LsFiles)
	patterns := pflag.Args()

	// For ls-files, if no patterns provided, just run the command
	// For track/untrack, patterns are required
	if err := lfsfiles.Execute(patterns, opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
