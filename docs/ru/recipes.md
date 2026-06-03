# Рецепты

Небольшие шелловые сниппеты, построенные на разделении чтения и записи.

## Все P1-задачи на сегодня

```sh
things today --json | jq '.[] | select(.tags | contains("P1"))'
```

## Массово навесить тег на список

```sh
things inbox --json | jq -r '.[].id' | while read id; do
  things tag "$id" "triage"
done
```

## Перенести всё с тегом `someday` в Someday

```sh
things anytime --json \
  | jq -r '.[] | select(.tags | contains("someday")) | .id' \
  | xargs -n1 -I{} things schedule {} someday
```

## Статусы проектов в таблицу

```sh
things projects --json \
  | jq -r '.[] | [.name, .area, .status] | @tsv' \
  | column -t -s$'\t'
```

## Аудит: открытые задачи без проекта и без области

```sh
things inbox --json \
  | jq '.[] | select(.project == "" and .area == "")'
```

## Быстрый capture из скрипта

```sh
new_todo() {
  things add "$1" --notes "$2" --tags "P3,inbox"
}
new_todo "Продлить домен" "namecheap"
```

## Предполётная проверка в автоматизации

Если встраиваешь это в персональный автомат, проваливай мягко, когда Things недоступен:

```sh
if ! things locale > /dev/null 2>&1; then
  echo "Things 3 недоступен — пропускаем capture"
  exit 0
fi
```
