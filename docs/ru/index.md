# thingscli

Тонкий CLI поверх AppleScript для [Things 3](https://culturedcode.com/things/). Команды чтения возвращают JSON, команды записи — короткий статус. Имена встроенных списков определяются рантайм-пробой, поэтому один и тот же бинарь работает с Things на любом языке интерфейса без отдельных флагов.

## Зачем

- **Пайплайны без регэкспов.** `things today --json | jq '.[] | select(.tags | contains("P1"))'` — никакого разбора локализованных строк.
- **Не зависим от `things:///` URL.** Полноценное чтение, фильтры, поиск — а не только write-and-forget.
- **Один бинарь, любая локаль.** Probe-then-cache: при первом вызове CLI спрашивает у Things, какие имена списков он реально использует; результат кэшируется в `~/.cache/thingscli/locale.json` на 30 дней.

## Кратко

```sh
things today                                 # JSON задач на сегодня
things search "ревью"                        # поиск по подстроке
things add "Оплатить аренду" --when tomorrow --tags "home"
things done <id>                             # пометить выполненной
things schedule <id> 2026-07-01              # запланировать
things --lang en inbox                       # принудительно фиксировать локаль
```

См. [Старт](getting-started.md) для установки и [обзор команд](commands/index.md).
