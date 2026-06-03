# Commands overview

Commands fall into two groups by output contract:

- **Read commands** emit JSON to stdout (array for lists, object for `show`/`locale`). They are safe to pipe into `jq`, scripts, or other tools. See [Reading](read.md).
- **Write commands** emit a single short human-readable status line on success (e.g. `DONE <id>`) and a non-zero exit on failure. See [Writing](write.md).

## Global flags

| Flag                | Purpose                                                                  |
|---------------------|--------------------------------------------------------------------------|
| `--lang <code>`     | Force a specific known language for the built-in list names.             |
| `--refresh-locale`  | Re-probe Things and rewrite the cached locale file.                      |
| `--pretty`          | Indent JSON output (read commands only).                                 |

## Identifying todos

All write commands take a Things to-do **id** as the first argument — the stable identifier returned by every read command as the `id` field. Never identify by name: names aren't unique and may change.

```sh
ID=$(things inbox | jq -r '.[0].id')
things done "$ID"
```

## Exit codes

- `0` — success.
- `1` — application error (Things returned an error, list not found, id not found, etc.). The message goes to stderr.
- `2` — usage error (wrong number of arguments, conflicting flags).
