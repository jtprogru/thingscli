# Recipes

Small composable shell snippets built on the read/write split.

## All P1 todos in Today

```sh
things today --json | jq '.[] | select(.tags | contains("P1"))'
```

## Bulk-tag a list

```sh
things inbox --json | jq -r '.[].id' | while read id; do
  things tag "$id" "triage"
done
```

## Move everything tagged `someday` into the Someday list

```sh
things anytime --json \
  | jq -r '.[] | select(.tags | contains("someday")) | .id' \
  | xargs -n1 -I{} things schedule {} someday
```

## Project status table

```sh
things projects --json \
  | jq -r '.[] | [.name, .area, .status] | @tsv' \
  | column -t -s$'\t'
```

## Audit: open todos with no project and no area

```sh
things inbox --json \
  | jq '.[] | select(.project == "" and .area == "")'
```

## Quick capture from a script

```sh
new_todo() {
  things add "$1" --notes "$2" --tags "P3,inbox"
}
new_todo "Renew domain" "namecheap"
```

## Pre-flight check in CI

If you wire this into a personal automation, fail fast when Things isn't reachable:

```sh
if ! things locale > /dev/null 2>&1; then
  echo "Things 3 not available — skipping personal capture"
  exit 0
fi
```
