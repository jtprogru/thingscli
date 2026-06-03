package cli

import (
	"encoding/json"
	"os"

	"github.com/jtprogru/thingscli/internal/things"
	"github.com/spf13/cobra"
)

var (
	flagLang    string
	flagRefresh bool
	flagPretty  bool
)

// NewRoot builds the root cobra command.
func NewRoot() *cobra.Command {
	root := &cobra.Command{
		Use:   "things",
		Short: "CLI for Things 3 via AppleScript (multi-locale).",
		Long: `things — thin CLI over AppleScript for Things 3.

Read commands print JSON so they parse cleanly.
Write commands print a short human-readable status line.

The built-in list names ("Inbox", "Today", ...) are auto-detected from
the running Things 3 app on first use, then cached. Override with --lang
or refresh with --refresh-locale.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentFlags().StringVar(&flagLang, "lang", "", "force locale (en, ru, de, fr, es, it, ja, zh-Hans, pt-BR, nl)")
	root.PersistentFlags().BoolVar(&flagRefresh, "refresh-locale", false, "ignore cached locale and re-probe Things")
	root.PersistentFlags().BoolVar(&flagPretty, "pretty", false, "pretty-print JSON output")

	addReadCommands(root)
	addWriteCommands(root)
	root.AddCommand(newVersionCmd())
	return root
}

func newClient() (*things.Client, error) {
	return things.New(flagLang, flagRefresh)
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	if flagPretty {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(v)
}
