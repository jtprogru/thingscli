# Reading

All commands here print JSON.

## Built-in lists

The seven Things built-in lists, addressed by canonical name regardless of UI language:

| Command            | Returns                                                              |
|--------------------|----------------------------------------------------------------------|
| `things inbox`     | Todos in Inbox.                                                      |
| `things today`     | Todos scheduled for today.                                           |
| `things upcoming`  | Todos scheduled later (the "Upcoming" view).                         |
| `things anytime`   | Todos in the Anytime view.                                           |
| `things someday`   | Todos in the Someday view.                                           |
| `things logbook`   | Completed/canceled archive.                                          |

Each yields a JSON array of `Todo` objects:

```json
[
  {
    "id": "5ZhX1arERehyn2JXGzd8Mk",
    "name": "Pay rent",
    "notes": "Use the second card",
    "status": "open",
    "tags": "home,P1",
    "due": "2026-07-01",
    "start": "",
    "project": "",
    "area": "Home"
  }
]
```

## Arbitrary lists

```sh
things list "<localized-list-name>"
```

Use this when you need a list whose name isn't one of the seven built-ins, or when you want to bypass the locale table entirely and address Things by its current display name.

## Projects, areas, tags

```sh
things projects     # [{id,name,area,status}, ...]
things areas        # [{id,name}, ...]
things tags         # ["P1","P2","home", ...]
```

## Project contents

```sh
things project "<Project Name>"
```

Returns todos belonging to the named project.

## Search

```sh
things search "<substring>"
```

Returns todos whose name contains the substring (case rules follow Things'). Searches across all lists.

## Single to-do

```sh
things show <id>
```

Returns one `Todo` object. Errors if no to-do has that id.

## Locale introspection

```sh
things locale
```

Returns the resolved locale table — useful for debugging when a built-in list command returns surprising results.

```json
{
  "lang": "ru",
  "inbox": "Входящие",
  "today": "Сегодня",
  "upcoming": "Завтра",
  "anytime": "В любое время",
  "someday": "Когда-нибудь",
  "logbook": "Журнал",
  "trash": "Корзина"
}
```
