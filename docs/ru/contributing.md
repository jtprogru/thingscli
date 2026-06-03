# Контрибьютинг

## Локальная разработка

```sh
git clone https://github.com/jtprogru/thingscli
cd thingscli
task build         # → ./dist/things
task ci            # lint + race-тесты (как в CI)
```

Нужны Go 1.26+, `go-task` и `golangci-lint v2.12+`. Для docs-сайта дополнительно — Python 3 и `task docs:install`.

## Добавить язык Things UI

Things 3 локализует имена семи встроенных списков. Чтобы добавить язык:

1. Открой Things на нужном языке UI и зафиксируй точные имена: **Входящие**, **Сегодня**, **Завтра/Upcoming**, **В любое время**, **Когда-нибудь**, **Журнал**, **Корзина**.
2. Добавь запись в `builtinLocales` в `internal/things/locale.go`. Поле `Lang` — BCP-47-style (`en`, `ru`, `pt-BR`, `zh-Hans`).
3. Прогон `things --lang <new-code> locale` подтвердит, что таблица грузится, а `things --refresh-locale locale` — что рантайм-проба её ловит на свежем кэше.

Всё. Отдельных кодовых веток на язык нет.

## Добавить команду

1. Реализуй операцию методом на `*things.Client` в `internal/things/ops.go`. Для всего, что возвращает задачи, переиспользуй `dumpListScript`.
2. Подключи к cobra в `internal/cli/read.go` или `internal/cli/write.go`. Чтение вызывает `printJSON`; запись печатает короткую строку `VERB <id>`.
3. Документируй в `docs/{en,ru}/commands/{read,write}.md` и, если неочевидно, добавь рецепт в `docs/{en,ru}/recipes.md`.

## Стиль

Репозиторий использует `golangci-lint v2`, конфиг — `.golangci.yaml`. Прогон `task lint` перед пушем. Из заметного:

- `gci` группирует импорты, с `github.com/jtprogru/thingscli` как отдельная секция.
- `gochecknoglobals` намеренно выключен для `internal/cli/*` (cobra-овые package-level vars) и `internal/things/*` (таблица локалей, константа AS-хелперов). Глобалки в других пакетах нужно обосновывать.
- Импорты идут блоками `stdlib | third-party | github.com/jtprogru/thingscli`, без инлайн- и префикс-комментариев.

## Коммиты

Conventional-префиксы (`feat:`, `fix:`, `chore(deps):`, `docs:`) держат GoReleaser-changelog аккуратным — релизные ноты по умолчанию фильтруют `docs:` и `test:`.

## Релизы

Тег от `main`:

```sh
git tag v0.X.Y
git push --tags
```

Дальше — workflow `goreleaser`: соберёт darwin/amd64 + darwin/arm64, подпишет чексаммы GPG-ключом проекта, обновит Homebrew-cask в `jtprogru/homebrew-tap`. Ручных шагов нет.
