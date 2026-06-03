# Architecture

## Layout

```
cmd/things/main.go      — entry point
internal/cli/           — cobra subcommands + flag wiring
internal/things/        — Things 3 client
  client.go             — osascript runner
  locale.go             — locale table + runtime probe + cache
  applescript.go        — shared AppleScript helpers + escape
  models.go             — Todo / Project / Area / ListKey
  ops.go                — read & write operations
```

The split keeps cobra-specific code out of the Things package, so the client is reusable as a library.

## Output protocol

AppleScript does not produce JSON. Embedding `jsonEscape` inside the script (the zsh original's approach) doubles every escape and is fragile. Instead, the AppleScript side emits records separated by ASCII control characters:

- `\x1e` (RS, ASCII 30) — between records
- `\x1f` (US, ASCII 31) — between fields within a record

These bytes do not occur in user-typed to-do content. Go reads the stream, splits on the separators, and marshals to JSON with `encoding/json`. Escaping happens once, on the Go side.

## Locale resolution

Things 3 names its built-in lists in the UI language. `list "Inbox"` works only when the UI is English; on a Russian install you need `list "Входящие"`. Hard-coding one set of names ties the tool to one language.

The CLI handles this in three steps:

1. **Lookup** — `internal/things/locale.go` defines `builtinLocales`, a table of `{lang, inbox, today, upcoming, ...}` records for the languages we know.
2. **Probe** — on first use, one `osascript` call iterates the candidate Inbox names and returns the first that Things actually resolves. That language becomes active.
3. **Cache** — the matched record is written to `~/.cache/thingscli/locale.json` with a timestamp. Subsequent invocations skip the probe for 30 days. `--refresh-locale` forces a re-probe; `--lang <code>` skips the cache and uses a specific table directly.

If Things uses a language we don't yet have, the probe returns empty and the user gets a clear error. Adding support is a one-line addition to `builtinLocales`.

## AppleScript embedding

Helpers (`pad`, `dateISO`, `parseDate`, `todoRecord`) live in one constant (`asHelpers`) that gets concatenated into each script needing them. Each script is otherwise a small, self-contained block of AppleScript built by Go string concatenation — no template engine, no nested escaping.

The only escaping we do is for AppleScript double-quoted literals: `\` → `\\`, `"` → `\"`. Everything else (multi-byte UTF-8, control characters that aren't `\x1e`/`\x1f`) passes through unchanged.

## Identity

Todos are addressed by their stable Things `id` (the same identifier the `things:///` URL scheme uses). Names are never used to identify — they aren't unique and may be edited.

## What this is not

- Not a `things:///` URL wrapper. The URL scheme is write-only and cannot drive reads/search/lookup; everything goes through AppleScript.
- Not a sync layer. There is no local database; every command is a live AppleScript call.
- Not destructive. There is no irreversible delete. `trash` is the strongest destructive op and it's reversible from inside Things.
