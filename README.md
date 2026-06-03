# thingscli

A thin Go CLI over AppleScript for [Things 3](https://culturedcode.com/things/). Read commands print a compact table by default (`--json` for the structured form used in pipelines); write commands emit a short status line. Built-in list names ("Inbox", "Today", …) are auto-detected by probing Things at runtime, so one binary works on any installed Things UI language.

[Documentation](https://jtprogru.github.io/thingscli/) · [Releases](https://github.com/jtprogru/thingscli/releases) · [Changelog](CHANGELOG.md)

## Install

```sh
brew install --cask jtprogru/tap/things
# or
go install github.com/jtprogru/thingscli/cmd/things@latest
```

## Use

```sh
things today                           # table of today's todos (id|name|status)
things today --json | jq               # raw JSON for pipelines
things search "review"                 # name substring search
things add "Pay rent" --when tomorrow --tags "home"
things done <id>                       # mark completed
things --lang ru inbox                 # force a known locale
```

Full reference: [Commands overview](https://jtprogru.github.io/thingscli/commands/).

## Multi-locale

Things 3 names its built-in lists in the UI language (`Inbox` / `Входящие` / `Eingang` / 受信箱 / …). The CLI handles that without per-locale flags:

1. **Probe** — first run asks Things, in one AppleScript call, which of the bundled Inbox names actually resolves. The match becomes the active language.
2. **Cache** — written to `~/.cache/thingscli/locale.json` for 30 days.
3. **Override** — `--lang <code>` skips the cache and uses a specific table; `--refresh-locale` re-probes.

Bundled languages: `en`, `ru`, `de`, `fr`, `es`, `it`, `ja`, `zh-Hans`, `pt-BR`, `nl`. Missing one? See [Contributing](https://jtprogru.github.io/thingscli/contributing/).

## Why not the `things:///` URL scheme?

The URL scheme is write-only. It can't list, search, or look up todos by id, which leaves it useless for any pipeline that needs to read state before deciding what to do. This tool uses AppleScript for both directions so the same identifier (`id`) flows from read to write without a translation layer.

## Develop

```sh
task build         # ./dist/things
task ci            # lint + race tests, matches CI
task release:dry   # local goreleaser snapshot in ./dist
task docs:serve    # MkDocs at http://127.0.0.1:8000
```

Go 1.26+, `go-task`, `golangci-lint v2.12+`. macOS only — the AppleScript bridge has no portable fallback.

## License

MIT — see [LICENSE](LICENSE).
