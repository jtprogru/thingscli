# Getting started

## Requirements

- macOS (Things 3 is mac-only; the CLI uses `/usr/bin/osascript`).
- Things 3 installed and able to launch.
- Go 1.26+ if building from source.

## Install

**Homebrew (recommended):**

```sh
brew install --cask jtprogru/tap/things
```

**From source:**

```sh
go install github.com/jtprogru/thingscli/cmd/things@latest
```

**From release archives:** download from the [GitHub Releases](https://github.com/jtprogru/thingscli/releases) page, extract, drop `things` into a directory on your `PATH`.

## First call

```sh
things locale
```

This triggers the runtime probe. The CLI asks Things, in a single AppleScript call, which of the known "Inbox" names actually resolves. The match becomes the active locale and is cached at `~/.cache/thingscli/locale.json` for 30 days.

If your Things UI language isn't in the bundled table, the probe will return an empty match — file an issue with the localized names of the 7 built-in lists and they will be added.

## Overrides

| Flag                | Effect                                                                  |
|---------------------|-------------------------------------------------------------------------|
| `--lang <code>`     | Skip auto-detect, use this language's name table (e.g. `--lang ru`).    |
| `--refresh-locale`  | Ignore the cache and re-probe Things.                                   |
| `--json`            | Emit raw JSON instead of the default table (read commands only).        |
| `--pretty`          | Indent JSON output. Only meaningful with `--json`.                      |

Known language codes: `en`, `ru`, `de`, `fr`, `es`, `it`, `ja`, `zh-Hans`, `pt-BR`, `nl`.

## Sanity check

```sh
things projects
things today
```

If both return a table (possibly empty), you're set. If `osascript` errors, make sure Things 3 is installed and has been opened at least once.
