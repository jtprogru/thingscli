package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/jtprogru/thingscli/internal/things"
)

const (
	idWidth   = 8
	nameWidth = 60
	ellipsis  = "…"
)

// shortID returns the first idWidth runes of a Things id. IDs are longer
// (typically 22 chars), but the prefix is unique enough to eyeball lists and
// to copy back into a write command via shell completion or paste.
func shortID(id string) string {
	r := []rune(id)
	if len(r) <= idWidth {
		return id
	}
	return string(r[:idWidth])
}

// truncName keeps a name within nameWidth runes, appending an ellipsis when
// it had to cut. Multi-byte handling: we count runes, not bytes, so Cyrillic
// and CJK names get the same visual budget as Latin ones.
func truncName(s string) string {
	r := []rune(s)
	if len(r) <= nameWidth {
		return s
	}
	return string(r[:nameWidth-1]) + ellipsis
}

func newTab(w io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
}

// printJSON writes v as JSON to stdout. Respects flagPretty.
func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	if flagPretty {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(v)
}

// renderTodos picks JSON or tabular output based on flagJSON.
func renderTodos(ts []things.Todo) error {
	if flagJSON {
		return printJSON(ts)
	}
	tw := newTab(os.Stdout)
	fmt.Fprintln(tw, "ID\tNAME\tSTATUS")
	for _, t := range ts {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", shortID(t.ID), truncName(t.Name), t.Status)
	}
	return tw.Flush()
}

// renderTodo prints a single to-do in JSON or as a key:value block.
func renderTodo(t things.Todo) error {
	if flagJSON {
		return printJSON(t)
	}
	tw := newTab(os.Stdout)
	rows := [][2]string{
		{"ID", t.ID},
		{"Name", t.Name},
		{"Status", t.Status},
		{"Tags", t.Tags},
		{"Due", t.Due},
		{"Start", t.Start},
		{"Project", t.Project},
		{"Area", t.Area},
		{"Notes", t.Notes},
	}
	for _, r := range rows {
		fmt.Fprintf(tw, "%s\t%s\n", r[0], r[1])
	}
	return tw.Flush()
}

func renderProjects(ps []things.Project) error {
	if flagJSON {
		return printJSON(ps)
	}
	tw := newTab(os.Stdout)
	fmt.Fprintln(tw, "ID\tNAME\tSTATUS")
	for _, p := range ps {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", shortID(p.ID), truncName(p.Name), p.Status)
	}
	return tw.Flush()
}

func renderAreas(as []things.Area) error {
	if flagJSON {
		return printJSON(as)
	}
	tw := newTab(os.Stdout)
	fmt.Fprintln(tw, "ID\tNAME")
	for _, a := range as {
		fmt.Fprintf(tw, "%s\t%s\n", shortID(a.ID), truncName(a.Name))
	}
	return tw.Flush()
}

func renderTags(tags []string) error {
	if flagJSON {
		return printJSON(tags)
	}
	for _, t := range tags {
		fmt.Fprintln(os.Stdout, t)
	}
	return nil
}

func renderLocale(l things.Locale) error {
	if flagJSON {
		return printJSON(l)
	}
	tw := newTab(os.Stdout)
	rows := [][2]string{
		{"Lang", l.Lang},
		{"Inbox", l.Inbox},
		{"Today", l.Today},
		{"Upcoming", l.Upcoming},
		{"Anytime", l.Anytime},
		{"Someday", l.Someday},
		{"Logbook", l.Logbook},
		{"Trash", l.Trash},
	}
	for _, r := range rows {
		fmt.Fprintf(tw, "%s\t%s\n", r[0], r[1])
	}
	return tw.Flush()
}
