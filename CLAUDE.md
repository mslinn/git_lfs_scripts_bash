# Claude Code Guidelines

## Flag Parsing

**POSIX compatibility is always required.**

All commands must use `github.com/spf13/pflag` instead of the standard library `flag` package to support POSIX-style flag combining.

### Example

```go
import (
    flag "github.com/spf13/pflag"
)

func main() {
    var dryRun, verbose bool

    // Use BoolVarP for shorthand and long form flags
    flag.BoolVarP(&dryRun, "dry-run", "d", false, "Dry run mode")
    flag.BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

    flag.Parse()
}
```

This allows users to combine flags: `-dv` instead of `-d -v`

## Dependencies

All direct dependencies of external tools (like giftless) must be checked before running the tool to provide clear error messages with installation instructions.
