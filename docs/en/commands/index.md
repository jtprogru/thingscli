# Commands overview

Commands fall into two groups by output contract:

- **Read commands** emit a compact table by default — three columns for to-do lists (`ID`, `NAME`, `STATUS`), key:value blocks for single objects (`show`, `locale`), one line per tag for `tags`. Pass `--json` for the raw JSON form used in pipelines. See [Reading](read.md).
- **Write commands** emit a single short human-readable status line on success (e.g. `DONE <id>`) and a non-zero exit on failure. The `--json` flag does not affect writes. See [Writing](write.md).

## Global flags

| Flag                | Purpose                                                                  |
|---------------------|--------------------------------------------------------------------------|
| `--lang <code>`     | Force a specific known language for the built-in list names.             |
| `--refresh-locale`  | Re-probe Things and rewrite the cached locale file.                      |
| `--json`            | Emit raw JSON instead of the default table (read commands only).         |
| `--pretty`          | Indent JSON output. Only meaningful with `--json`.                       |

## Identifying todos

All write commands take a Things to-do **id** as the first argument — the stable identifier. The default table view shows only the first 8 characters of each id for readability; the full id is available via `--json` or via `things show <prefix>`. Never identify by name: names aren't unique and may change.

```sh
ID=$(things inbox --json | jq -r '.[0].id')
things done "$ID"
```

## Exit codes

- `0` — success.
- `1` — application error (Things returned an error, list not found, id not found, etc.). The message goes to stderr.
- `2` — usage error (wrong number of arguments, conflicting flags).
