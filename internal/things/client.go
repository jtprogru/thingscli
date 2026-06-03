package things

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// runOsa executes an AppleScript via osascript and returns trimmed stdout.
// stderr is included in the error on failure so callers see Things' own
// error messages (e.g., "Can't get list \"Foo\"").
func runOsa(script string) (string, error) {
	cmd := exec.CommandContext(context.Background(), "/usr/bin/osascript", "-l", "AppleScript", "-e", script)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			return "", fmt.Errorf("osascript: %w", err)
		}
		return "", fmt.Errorf("osascript: %s", msg)
	}
	return strings.TrimRight(stdout.String(), "\n"), nil
}

// Client is the high-level interface to Things 3.
type Client struct {
	Locale Locale
}

// New returns a Client with the active Things locale resolved.
// langOverride forces a specific known language; pass "" for auto-detect.
// forceRefresh bypasses the locale cache.
func New(langOverride string, forceRefresh bool) (*Client, error) {
	loc, err := ResolveLocale(langOverride, forceRefresh)
	if err != nil {
		return nil, err
	}
	return &Client{Locale: loc}, nil
}
