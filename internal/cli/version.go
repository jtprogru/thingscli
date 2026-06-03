package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Build-time metadata. Overridden by goreleaser's -ldflags.
// Keep these as package-level vars (not consts) so -X can set them.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
	BuiltBy = "source"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version, commit, build date, and Go runtime info",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("things %s\n", Version)
			fmt.Printf("  commit:    %s\n", Commit)
			fmt.Printf("  built:     %s by %s\n", Date, BuiltBy)
			fmt.Printf("  go:        %s\n", runtime.Version())
			fmt.Printf("  platform:  %s/%s\n", runtime.GOOS, runtime.GOARCH)
		},
	}
}
