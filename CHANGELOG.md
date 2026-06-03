# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

This file is the project-level human-curated history. The GoReleaser pipeline still auto-generates per-release notes on the GitHub Releases page from commit messages.

## [Unreleased]

### Changed

- **Default output for read commands is now a compact table** instead of JSON. Lists (`today`, `inbox`, `upcoming`, `anytime`, `someday`, `logbook`, `list`, `project`, `search`) print three columns: `ID` (first 8 chars of the Things id), `NAME` (truncated at 60 runes with `…`), and `STATUS`. `show` and `locale` print key:value blocks. `projects` and `areas` print their own tables; `tags` prints one tag per line.
- New `--json` flag restores the previous behavior. Combine with `--pretty` for indented JSON. Existing pipelines should add `--json` to their `things …` invocations.
- `--pretty` now only takes effect with `--json` (no functional change for callers who already paired them, but the help text reflects the dependency).

## [0.1.0] - 2026-06-04

### Added

- Initial Go implementation of the CLI, ported from the prior zsh prototype.
- Read commands: `today`, `inbox`, `upcoming`, `anytime`, `someday`, `logbook`, `list`, `project`, `search`, `show`, `projects`, `areas`, `tags`, `locale`.
- Write commands: `add`, `done`, `cancel`, `reopen`, `rename`, `note`, `move`, `schedule`, `due`, `tag`, `trash`.
- Runtime locale probe: detects the active Things UI language by asking Things which built-in list names resolve, caches the result at `~/.cache/thingscli/locale.json` for 30 days. `--lang` overrides; `--refresh-locale` re-probes.
- Built-in locale tables for `en`, `ru`, `de`, `fr`, `es`, `it`, `ja`, `zh-Hans`, `pt-BR`, `nl`.
- `version` command exposing build-time metadata injected by GoReleaser via `-ldflags`.
- `Taskfile.yml` with `run`, `build`, `lint`, `test`, `test:race`, `ci`, `release:dry`, and `docs:*` targets.
- GoReleaser config restricted to `darwin/{amd64,arm64}` (the AppleScript bridge is mac-only by construction), GPG-signed checksums, and Homebrew cask publishing to `jtprogru/homebrew-tap`.
- `golangci-lint v2` config aligned with the project's sibling `srekit`.
- GitHub Actions workflows: lint, test (race + coverage), GoReleaser on tag, MkDocs deploy.
- Dependabot for `gomod` and `github-actions`, weekly cadence.
- MkDocs Material site with `mkdocs-static-i18n` (English + Russian).
