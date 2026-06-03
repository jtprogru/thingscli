# Writing

All commands here print one short status line on success and exit non-zero on failure.

## Create

```sh
things add "<title>" [flags]
```

| Flag         | Effect                                                                             |
|--------------|------------------------------------------------------------------------------------|
| `--notes`    | Set the to-do's notes.                                                             |
| `--tags`     | Comma-separated tag list (tags must already exist in Things).                      |
| `--project`  | Place the new to-do in this project.                                               |
| `--area`     | Place it in this area.                                                             |
| `--list`     | Place it in this built-in list, addressed by its localized name.                   |
| `--when`     | Schedule: `today`, `tomorrow`, or `YYYY-MM-DD`.                                    |

`--project`, `--area`, `--list` are mutually exclusive.

Prints `ADDED <id> :: <title>`.

## Status changes

```sh
things done   <id>
things cancel <id>
things reopen <id>
```

Map to the AppleScript `status` field (`completed`, `canceled`, `open`).

## Edit

```sh
things rename <id> "<new title>"
things note   <id> "<new notes>"      # overwrite (not append)
things tag    <id> "tag1,tag2"        # overwrite tag list
```

## Move

```sh
things move <id> --project "<name>"
things move <id> --area    "<name>"
things move <id> --list    "<name>"
```

Exactly one of `--project|--area|--list` must be set.

## Schedule

```sh
things schedule <id> today
things schedule <id> tomorrow
things schedule <id> 2026-07-01
things schedule <id> someday
```

`someday` moves the to-do into the Someday list (using the resolved locale name).

## Due date

```sh
things due <id> 2026-07-15
things due <id> clear
```

## Delete (reversible)

```sh
things trash <id>
```

Moves to Things' Trash list. Recover from inside Things if you change your mind. There is intentionally no `empty-trash` and no irreversible delete.
