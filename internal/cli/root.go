package cli

import (
	"github.com/jtprogru/thingscli/internal/things"
	"github.com/spf13/cobra"
)

var (
	flagLang    string
	flagRefresh bool
	flagJSON    bool
	flagPretty  bool
)

// NewRoot builds the root cobra command.
func NewRoot() *cobra.Command {
	root := &cobra.Command{
		Use:   "things",
		Short: "CLI for Things 3 via AppleScript (multi-locale).",
		Long: `things — thin CLI over AppleScript for Things 3.

Read commands print a compact table by default. Pass --json for the
raw JSON form suitable for piping into jq or other tools; combine
with --pretty for indented JSON. Write commands print a short
human-readable status line.

The built-in list names ("Inbox", "Today", ...) are auto-detected from
the running Things 3 app on first use, then cached. Override with --lang
or refresh with --refresh-locale.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentFlags().StringVar(&flagLang, "lang", "", "force locale (en, ru, de, fr, es, it, ja, zh-Hans, pt-BR, nl)")
	root.PersistentFlags().BoolVar(&flagRefresh, "refresh-locale", false, "ignore cached locale and re-probe Things")
	root.PersistentFlags().BoolVar(&flagJSON, "json", false, "emit raw JSON instead of a table (read commands only)")
	root.PersistentFlags().BoolVar(&flagPretty, "pretty", false, "indent JSON output (only meaningful with --json)")

	addReadCommands(root)
	addWriteCommands(root)
	root.AddCommand(newVersionCmd())
	return root
}

func newClient() (*things.Client, error) {
	return things.New(flagLang, flagRefresh)
}
