# Reading

By default, commands here print a compact table. Pass `--json` for the structured form (a JSON array for lists, an object for `show`/`locale`, a flat list of strings for `tags`). Combine `--json` with `--pretty` for indented JSON.

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

Default output is a three-column table:

```
ID        NAME                          STATUS
5ZhX1arE  Pay rent                      open
WavUHfP2  Read the Linux kernel post    open
```

The `ID` column shows the first 8 characters of the stable Things id — enough to eyeball but not the full identifier. Pass `--json` to get the full record:

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
things projects     # table: ID | NAME | STATUS
things areas        # table: ID | NAME
things tags         # one tag per line
```

Add `--json` for `[{id,name,area,status}, ...]`, `[{id,name}, ...]`, and `["P1","P2",...]` respectively.

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

Default output is a key:value block with every field:

```
ID       5ZhX1arERehyn2JXGzd8Mk
Name     Pay rent
Status   open
Tags     home,P1
Due      2026-07-01
Start
Project
Area     Home
Notes    Use the second card
```

Pass `--json` for one `Todo` object. Errors if no to-do has that id.

## Locale introspection

```sh
things locale
```

The resolved locale table — useful for debugging when a built-in list command returns surprising results.

```
Lang      ru
Inbox     Входящие
Today     Сегодня
Upcoming  Завтра
Anytime   В любое время
Someday   Когда-нибудь
Logbook   Журнал
Trash     Корзина
```

Add `--json` for the structured form.
