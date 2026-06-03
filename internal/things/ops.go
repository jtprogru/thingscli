package things

import (
	"fmt"
	"strings"
)

// --- read ------------------------------------------------------------------

// dumpListScript wraps an arbitrary AppleScript body that must assign
// `theTodos` to a list of to dos, and emits a recordSep-terminated stream
// of fieldSep-joined to-do records.
func dumpListScript(body string) string {
	return asHelpers + `
set RS to character id 30
set out to ""
tell application "Things3"
` + body + `
end tell
repeat with t in theTodos
	set out to out & (my todoRecord(t)) & RS
end repeat
return out
`
}

func parseTodos(out string) []Todo {
	if out == "" {
		return []Todo{}
	}
	records := strings.Split(out, recordSep)
	todos := make([]Todo, 0, len(records))
	for _, rec := range records {
		if rec == "" {
			continue
		}
		f := strings.Split(rec, fieldSep)
		// Defensive: must have 9 fields.
		for len(f) < 9 {
			f = append(f, "")
		}
		todos = append(todos, Todo{
			ID: f[0], Name: f[1], Notes: f[2], Status: f[3],
			Tags: f[4], Due: f[5], Start: f[6],
			Project: f[7], Area: f[8],
		})
	}
	return todos
}

func (c *Client) listByName(name string) ([]Todo, error) {
	body := `	set theTodos to to dos of list "` + escapeAS(name) + `"`
	out, err := runOsa(dumpListScript(body))
	if err != nil {
		return nil, err
	}
	return parseTodos(out), nil
}

// BuiltinList returns todos from one of the canonical built-in lists,
// using the active locale's name for it.
func (c *Client) BuiltinList(k ListKey) ([]Todo, error) {
	name, err := c.Locale.Name(k)
	if err != nil {
		return nil, err
	}
	return c.listByName(name)
}

// List returns todos from any list by its exact (localized) name.
func (c *Client) List(name string) ([]Todo, error) {
	return c.listByName(name)
}

// Project returns todos of the given project name.
func (c *Client) Project(name string) ([]Todo, error) {
	body := `	set theTodos to to dos of project "` + escapeAS(name) + `"`
	out, err := runOsa(dumpListScript(body))
	if err != nil {
		return nil, err
	}
	return parseTodos(out), nil
}

// Search returns todos whose name contains the substring.
func (c *Client) Search(q string) ([]Todo, error) {
	body := `	set theTodos to (to dos whose name contains "` + escapeAS(q) + `")`
	out, err := runOsa(dumpListScript(body))
	if err != nil {
		return nil, err
	}
	return parseTodos(out), nil
}

// Show returns a single to-do by ID.
func (c *Client) Show(id string) (Todo, error) {
	script := asHelpers + `
set RS to character id 30
tell application "Things3"
	set t to to do id "` + escapeAS(id) + `"
end tell
return (my todoRecord(t)) & RS
`
	out, err := runOsa(script)
	if err != nil {
		return Todo{}, err
	}
	todos := parseTodos(out)
	if len(todos) == 0 {
		return Todo{}, fmt.Errorf("no to-do with id %q", id)
	}
	return todos[0], nil
}

// Projects returns all projects.
func (c *Client) Projects() ([]Project, error) {
	script := `
set FS to character id 31
set RS to character id 30
set out to ""
tell application "Things3"
	repeat with p in projects
		set areaName to ""
		try
			set areaName to name of area of p
		end try
		set out to out & (id of p) & FS & (name of p) & FS & areaName & FS & ((status of p) as string) & RS
	end repeat
end tell
return out
`
	raw, err := runOsa(script)
	if err != nil {
		return nil, err
	}
	if raw == "" {
		return []Project{}, nil
	}
	records := strings.Split(raw, recordSep)
	out := make([]Project, 0, len(records))
	for _, rec := range records {
		if rec == "" {
			continue
		}
		f := strings.Split(rec, fieldSep)
		for len(f) < 4 {
			f = append(f, "")
		}
		out = append(out, Project{ID: f[0], Name: f[1], Area: f[2], Status: f[3]})
	}
	return out, nil
}

// Areas returns all areas.
func (c *Client) Areas() ([]Area, error) {
	script := `
set FS to character id 31
set RS to character id 30
set out to ""
tell application "Things3"
	repeat with a in areas
		set out to out & (id of a) & FS & (name of a) & RS
	end repeat
end tell
return out
`
	raw, err := runOsa(script)
	if err != nil {
		return nil, err
	}
	if raw == "" {
		return []Area{}, nil
	}
	records := strings.Split(raw, recordSep)
	out := make([]Area, 0, len(records))
	for _, rec := range records {
		if rec == "" {
			continue
		}
		f := strings.Split(rec, fieldSep)
		for len(f) < 2 {
			f = append(f, "")
		}
		out = append(out, Area{ID: f[0], Name: f[1]})
	}
	return out, nil
}

// Tags returns all tag names.
func (c *Client) Tags() ([]string, error) {
	script := `
set RS to character id 30
set out to ""
tell application "Things3"
	repeat with tg in tags
		set out to out & (name of tg) & RS
	end repeat
end tell
return out
`
	raw, err := runOsa(script)
	if err != nil {
		return nil, err
	}
	if raw == "" {
		return []string{}, nil
	}
	tags := []string{}
	for _, t := range strings.Split(raw, recordSep) {
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags, nil
}

// --- write -----------------------------------------------------------------

// AddOptions configures a new to-do.
type AddOptions struct {
	Notes   string
	Project string
	Area    string
	List    string // built-in list name (localized)
	When    string // "today" | "tomorrow" | "YYYY-MM-DD" | ""
	Tags    string // comma-separated
}

// Add creates a new to-do. Returns the new ID.
func (c *Client) Add(title string, opt AddOptions) (string, error) {
	props := `name:"` + escapeAS(title) + `"`
	if opt.Notes != "" {
		props += `, notes:"` + escapeAS(opt.Notes) + `"`
	}
	if opt.Tags != "" {
		props += `, tag names:"` + escapeAS(opt.Tags) + `"`
	}

	var mover string
	switch {
	case opt.Project != "":
		mover = `	set project of newTodo to project "` + escapeAS(opt.Project) + `"`
	case opt.Area != "":
		mover = `	move newTodo to area "` + escapeAS(opt.Area) + `"`
	case opt.List != "":
		mover = `	move newTodo to list "` + escapeAS(opt.List) + `"`
	}

	var scheduler string
	switch opt.When {
	case "":
		// nothing
	case "today":
		scheduler = `	schedule newTodo for (current date)`
	case "tomorrow":
		scheduler = `	schedule newTodo for ((current date) + 1 * days)`
	default:
		scheduler = `	schedule newTodo for (my parseDate("` + escapeAS(opt.When) + `"))`
	}

	script := asHelpers + `
tell application "Things3"
	set newTodo to make new to do with properties {` + props + `}
` + mover + `
` + scheduler + `
	return id of newTodo
end tell
`
	return runOsa(script)
}

func (c *Client) setStatus(id, status string) error {
	script := `tell application "Things3" to set status of to do id "` + escapeAS(id) + `" to ` + status
	_, err := runOsa(script)
	return err
}

func (c *Client) Done(id string) error   { return c.setStatus(id, "completed") }
func (c *Client) Cancel(id string) error { return c.setStatus(id, "canceled") }
func (c *Client) Reopen(id string) error { return c.setStatus(id, "open") }

func (c *Client) Rename(id, newName string) error {
	script := `tell application "Things3" to set name of to do id "` + escapeAS(id) + `" to "` + escapeAS(newName) + `"`
	_, err := runOsa(script)
	return err
}

func (c *Client) Note(id, text string) error {
	script := `tell application "Things3" to set notes of to do id "` + escapeAS(id) + `" to "` + escapeAS(text) + `"`
	_, err := runOsa(script)
	return err
}

// MoveTarget identifies where to move a to-do.
type MoveTarget struct {
	Kind string // "project" | "area" | "list"
	Name string
}

// Move relocates a to-do.
func (c *Client) Move(id string, t MoveTarget) error {
	eid := escapeAS(id)
	etarget := escapeAS(t.Name)
	var script string
	if t.Kind == "project" {
		script = `tell application "Things3" to set project of to do id "` + eid + `" to project "` + etarget + `"`
	} else {
		script = `tell application "Things3" to move (to do id "` + eid + `") to ` + t.Kind + ` "` + etarget + `"`
	}
	_, err := runOsa(script)
	return err
}

// Schedule sets activation date. when = "today" | "tomorrow" | "YYYY-MM-DD" | "someday".
func (c *Client) Schedule(id, when string) error {
	eid := escapeAS(id)
	var expr string
	switch when {
	case "today":
		expr = `schedule (to do id "` + eid + `") for (current date)`
	case "tomorrow":
		expr = `schedule (to do id "` + eid + `") for ((current date) + 1 * days)`
	case "someday":
		name, err := c.Locale.Name(ListSomeday)
		if err != nil {
			return err
		}
		expr = `move (to do id "` + eid + `") to list "` + escapeAS(name) + `"`
	default:
		expr = `schedule (to do id "` + eid + `") for (my parseDate("` + escapeAS(when) + `"))`
	}
	script := asHelpers + `
tell application "Things3"
	` + expr + `
end tell
`
	_, err := runOsa(script)
	return err
}

// Due sets or clears the due date. Pass "clear" to remove.
func (c *Client) Due(id, date string) error {
	eid := escapeAS(id)
	if date == "clear" {
		script := `tell application "Things3" to set due date of to do id "` + eid + `" to missing value`
		_, err := runOsa(script)
		return err
	}
	script := asHelpers + `
tell application "Things3" to set due date of to do id "` + eid + `" to (my parseDate("` + escapeAS(date) + `"))
`
	_, err := runOsa(script)
	return err
}

// Tag replaces the comma-separated tag list on a to-do.
func (c *Client) Tag(id, tags string) error {
	script := `tell application "Things3" to set tag names of to do id "` + escapeAS(id) + `" to "` + escapeAS(tags) + `"`
	_, err := runOsa(script)
	return err
}

// Trash moves a to-do to Things' trash list (reversible).
func (c *Client) Trash(id string) error {
	name, err := c.Locale.Name(ListTrash)
	if err != nil {
		return err
	}
	script := `tell application "Things3" to move (to do id "` + escapeAS(id) + `") to list "` + escapeAS(name) + `"`
	_, err = runOsa(script)
	return err
}
