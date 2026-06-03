# Contributing

## Dev setup

```sh
git clone https://github.com/jtprogru/thingscli
cd thingscli
task build         # → ./dist/things
task ci            # lint + race tests (matches CI)
```

Go 1.26+, `go-task`, and `golangci-lint v2.12+` are expected. For the docs site you'll also need Python 3 and `task docs:install`.

## Adding a Things UI language

Things 3 localizes its seven built-in list names. To add a language:

1. Open Things in that UI language and note the exact names of: **Inbox**, **Today**, **Upcoming**, **Anytime**, **Someday**, **Logbook**, **Trash**.
2. Append a record to `builtinLocales` in `internal/things/locale.go`. The `Lang` field uses BCP-47-ish codes (`en`, `ru`, `pt-BR`, `zh-Hans`).
3. Run `things --lang <new-code> locale` to verify the table loads, then `things --refresh-locale locale` to confirm the runtime probe picks it up on a fresh cache.

That's the whole change — there is no per-language code path.

## Adding a command

1. Implement the operation as a method on `*things.Client` in `internal/things/ops.go`. Reuse `dumpListScript` for anything returning todos.
2. Wire it into cobra in `internal/cli/read.go` or `internal/cli/write.go`. Read commands call `printJSON`; write commands print a short `VERB <id>` status line.
3. Document it in `docs/en/commands/{read,write}.md` and (if non-obvious) add a recipe in `docs/en/recipes.md`.

## Style

The repo uses `golangci-lint v2` with the config in `.golangci.yaml`. Run `task lint` before pushing. Notable rules:

- `gci` enforces import grouping with `github.com/jtprogru/thingscli` as a separate section.
- `gochecknoglobals` is intentionally off for `internal/cli/*` (cobra package-level vars) and `internal/things/*` (the locale table, the AppleScript helpers constant). New globals in other packages need justification.
- Imports must group `stdlib | third-party | github.com/jtprogru/thingscli` with no inline or prefix comments.

## Commits

Conventional-style prefixes (`feat:`, `fix:`, `chore(deps):`, `docs:`) keep the GoReleaser changelog tidy — the release notes filter out `docs:` and `test:` by design.

## Releases

Tag from `main`:

```sh
git tag v0.X.Y
git push --tags
```

The `goreleaser` workflow handles the rest: builds darwin/amd64 + darwin/arm64, signs checksums with the project GPG key, updates the Homebrew cask in `jtprogru/homebrew-tap`. There is no manual step.
