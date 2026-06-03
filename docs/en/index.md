# thingscli

A thin CLI over AppleScript for [Things 3](https://culturedcode.com/things/). Read commands return JSON for safe piping; write commands print a short status line. Built-in list names are auto-detected by probing Things at runtime, so the same binary works on any installed Things UI language without per-locale flags.

## Why

- **Pipelines without regex.** `things today | jq '.[] | select(.tags | contains("P1"))'` — no parsing localized strings.
- **No reliance on `things:///` URLs.** Full read API for filters/search/lookup, not only fire-and-forget writes.
- **Single binary, multi-locale.** Probe-then-cache: the first call asks Things which list names it actually uses; the result lives in `~/.cache/thingscli/locale.json` for 30 days.

## At a glance

```sh
things today                                 # JSON of today's todos
things search "review"                       # name substring search
things add "Pay rent" --when tomorrow --tags "home"
things done <id>                             # mark completed
things schedule <id> 2026-07-01              # date-pick
things --lang de inbox                       # force a known locale
```

See [Getting started](getting-started.md) for installation and the [Commands overview](commands/index.md) for the full surface.
