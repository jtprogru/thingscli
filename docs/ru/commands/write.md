# Запись

Все команды печатают одну короткую строку статуса при успехе и возвращают ненулевой exit на ошибке.

## Создать

```sh
things add "<заголовок>" [флаги]
```

| Флаг         | Эффект                                                                                |
|--------------|----------------------------------------------------------------------------------------|
| `--notes`    | Заметка.                                                                               |
| `--tags`     | Теги через запятую (должны уже существовать в Things).                                 |
| `--project`  | Положить в проект.                                                                     |
| `--area`     | Положить в область.                                                                    |
| `--list`     | Положить во встроенный список (по локализованному имени).                              |
| `--when`     | Запланировать: `today`, `tomorrow` или `YYYY-MM-DD`.                                   |

`--project`, `--area`, `--list` взаимоисключающие.

Печатает `ADDED <id> :: <заголовок>`.

## Смена статуса

```sh
things done   <id>
things cancel <id>
things reopen <id>
```

Маппятся в AppleScript-поле `status` (`completed`, `canceled`, `open`).

## Правка

```sh
things rename <id> "<новый заголовок>"
things note   <id> "<новый текст>"      # перезаписывает, не добавляет
things tag    <id> "tag1,tag2"          # перезаписывает список тегов
```

## Перемещение

```sh
things move <id> --project "<имя>"
things move <id> --area    "<имя>"
things move <id> --list    "<имя>"
```

Ровно один из `--project|--area|--list` должен быть задан.

## Планирование

```sh
things schedule <id> today
things schedule <id> tomorrow
things schedule <id> 2026-07-01
things schedule <id> someday
```

`someday` перемещает задачу в список «Когда-нибудь» (по резолвленному имени локали).

## Срок (due date)

```sh
things due <id> 2026-07-15
things due <id> clear
```

## Удаление (обратимое)

```sh
things trash <id>
```

Перемещает в Корзину Things. Восстановить можно из самого Things. Необратимого удаления нет принципиально — нет `empty-trash`, нет hard delete.
