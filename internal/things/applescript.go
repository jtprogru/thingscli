package things

import (
	"strings"
)

// AppleScript control characters used as field/record separators in script
// output. They are guaranteed not to appear in user-typed to-do content.
const (
	recordSep = "\x1e" // ASCII 30
	fieldSep  = "\x1f" // ASCII 31
)

// escapeAS escapes a Go string so it can be embedded inside an AppleScript
// double-quoted literal.
func escapeAS(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

// asHelpers is a block of AppleScript helper handlers reused across scripts.
// It defines:
//   - pad(n)      : zero-pad single-digit integers to two characters.
//   - dateISO(d)  : convert an AppleScript date to "YYYY-MM-DD", empty if missing.
//   - parseDate(s): parse "YYYY-MM-DD" into an AppleScript date at 00:00:00.
//   - todoRecord(t): serialize a to-do as a fieldSep-joined string of nine fields.
const asHelpers = `
on pad(n)
	set s to n as string
	if (count of s) = 1 then return "0" & s
	return s
end pad

on dateISO(d)
	if d is missing value then return ""
	return (year of d as string) & "-" & my pad(month of d as integer) & "-" & my pad(day of d as integer)
end dateISO

on parseDate(s)
	set y to (text 1 thru 4 of s) as integer
	set m to (text 6 thru 7 of s) as integer
	set d to (text 9 thru 10 of s) as integer
	set dt to current date
	set time of dt to 0
	set year of dt to y
	set month of dt to m
	set day of dt to d
	return dt
end parseDate

on todoRecord(t)
	set FS to character id 31
	tell application "Things3"
		set theId to id of t
		set theName to name of t
		set theNotes to notes of t
		set theStatus to (status of t) as string
		set theTags to tag names of t
		set theDue to my dateISO(due date of t)
		set theStart to my dateISO(activation date of t)
		set theProj to ""
		try
			set theProj to name of project of t
		end try
		set theArea to ""
		try
			set theArea to name of area of t
		end try
	end tell
	return theId & FS & theName & FS & theNotes & FS & theStatus & FS & theTags & FS & theDue & FS & theStart & FS & theProj & FS & theArea
end todoRecord
`

// buildProbeScript returns an AppleScript that asks Things 3 which of the
// given list names exists, returning the first match. We try them in order
// so the explicit order in builtinLocales (English first) is the tie-breaker.
func buildProbeScript(names []string) string {
	var sb strings.Builder
	sb.WriteString(`tell application "Things3"` + "\n")
	for _, n := range names {
		sb.WriteString(`	try` + "\n")
		sb.WriteString(`		get list "` + escapeAS(n) + `"` + "\n")
		sb.WriteString(`		return "` + escapeAS(n) + `"` + "\n")
		sb.WriteString(`	end try` + "\n")
	}
	sb.WriteString(`end tell` + "\n")
	sb.WriteString(`return ""` + "\n")
	return sb.String()
}
