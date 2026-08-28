# named-resource-service — curl examples

These examples assume the server is running locally on port `8080`:

```bash
rm -r data; mkdir -p data
./bin/named-resource-service serve --port 8080 --database-url "file:./data/named-resource.db?_fk=1"
```

Set a base URL for convenience:

```bash
BASE=http://localhost:8080
```

The `Named` resource lives at the `/nameds` collection. Its spec fields are:

| Field                | Type   | Notes                                  |
| -------------------- | ------ | -------------------------------------- |
| `name`               | string | **unique**                             |
| `altName`            | string | not null, max size 64                  |
| `number`             | int    | defaults to 0                          |
| `somethingOrNothing` | string | nullable, optional                     |

Resources are addressed by their generated `uid` (found in `metadata.uid`),
not by `name`.

---

## POST — create a resource

```bash
curl -s -X POST "$BASE/nameds" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "widget-a" },
    "spec": {
      "name": "widget-a",
      "altName": "Widget A",
      "number": 1,
      "somethingOrNothing": "hello"
    }
  }' | jq
```

The response includes the generated `metadata.uid`. Capture it for the
follow-up requests:

```bash
RESOURCE_UID=$(curl -s -X POST "$BASE/nameds" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "widget-b" },
    "spec": {
      "name": "widget-b",
      "altName": "Widget B",
      "number": 2
    }
  }' | jq -r '.metadata.uid')

echo "Created RESOURCE_UID: $RESOURCE_UID"
```

---

## GET — list all resources

```bash
curl -s "$BASE/nameds" | jq
```

## GET — fetch a single resource by UID

```bash
curl -s "$BASE/nameds/$RESOURCE_UID" | jq
```

---

## PUT — update a resource

Replace the resource identified by `$RESOURCE_UID`. Include the full desired spec.

```bash
curl -s -X PUT "$BASE/nameds/$RESOURCE_UID" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "widget-b" },
    "spec": {
      "name": "widget-b",
      "altName": "Widget B (updated)",
      "number": 42,
      "somethingOrNothing": "updated value"
    }
  }' | jq
```

---

## DELETE — remove a resource

```bash
curl -s -X DELETE "$BASE/nameds/$RESOURCE_UID" | jq
```

If the response body is empty on delete, you can confirm removal with:

```bash
curl -s -o /dev/null -w '%{http_code}\n' "$BASE/nameds/$RESOURCE_UID"
```

---

## Name uniqueness enforcement test

The `name` spec field is declared unique. Posting two resources with the same
`spec.name` must fail on the second attempt.

```bash
# First create with name "dupe" — expected to succeed.
echo "First POST (should succeed):"
curl -s -X POST "$BASE/nameds" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "dupe" },
    "spec": {
      "name": "dupe",
      "altName": "Duplicate Test",
      "number": 1
    }
  }' | jq

# Second create with the same spec.name "dupe" — expected to FAIL.
echo "Second POST (should fail on uniqueness):"
curl -s -X POST "$BASE/nameds" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "dupe-2" },
    "spec": {
      "name": "dupe",
      "altName": "Duplicate Test 2",
      "number": 2
    }
  }' | jq
```

To make the failure explicit, capture the HTTP status codes:

```bash
echo "First POST status:"
curl -s -o /dev/null -w '%{http_code}\n' -X POST "$BASE/nameds" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "unique-check" },
    "spec": { "name": "unique-check", "altName": "First", "number": 1 }
  }'

echo "Second POST status (expect a 4xx/5xx, not 2xx):"
curl -s -o /dev/null -w '%{http_code}\n' -X POST "$BASE/nameds" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "unique-check-again" },
    "spec": { "name": "unique-check", "altName": "Second", "number": 2 }
  }'
```

The first POST returns a `2xx` status; the second returns an error status
because the unique constraint on `spec.name` rejects the duplicate.

---

## Inspecting the SQLite database

The SQLite backend writes to the file named in `--database-url`
(`data/named-resource.db` above). Inspect it with the `sqlite3` CLI.

No local `sqlite3` install? Run it from a container, mounting the project dir
so the `.db` file is visible. `alpine` is available everywhere (install sqlite
on the fly); `keinos/sqlite3` is a dedicated image with `sqlite3` as its
entrypoint:

```bash
# alpine (ubiquitous): install sqlite, then run
docker run --rm -it -v "$PWD:/work" -w /work alpine \
  sh -c 'apk add --no-cache sqlite >/dev/null && sqlite3 data/named-resource.db'

# one-off query with alpine
docker run --rm -v "$PWD:/work" -w /work alpine \
  sh -c 'apk add --no-cache sqlite >/dev/null && sqlite3 data/named-resource.db ".tables"'

# dedicated image (sqlite3 is the entrypoint — pass args directly)
docker run --rm -it -v "$PWD:/work" -w /work keinos/sqlite3 data/named-resource.db
docker run --rm    -v "$PWD:/work" -w /work keinos/sqlite3 data/named-resource.db '.tables'
```

To avoid reinstalling sqlite on every `alpine` run, build a reusable image once
(one-liner, Dockerfile piped over stdin) and reference it afterwards:

```bash
printf 'FROM alpine\nRUN apk add --no-cache sqlite\n' | docker build -t alpine-sqlite -

# then use it directly — no apk step needed
docker run --rm -it -v "$PWD:/work" -w /work alpine-sqlite sqlite3 data/named-resource.db
docker run --rm    -v "$PWD:/work" -w /work alpine-sqlite sqlite3 data/named-resource.db '.tables'
```

The examples below use a local `sqlite3`; prefix any of them with one of the
`docker run ...` wrappers above to run them from a container instead.

Open an interactive shell:

```bash
sqlite3 data/named-resource.db
```

Useful dot-commands inside the shell:

```
.tables            -- list tables (resources, labels, annotations)
.schema            -- full schema for every table
.schema resources  -- schema for one table
.indexes resources -- indexes on a table
.headers on        -- show column names in query output
.mode column       -- align output into columns
.quit              -- exit
```

Run one-off queries without entering the shell by passing the SQL (or a
dot-command) as an argument:

```bash
# List tables
sqlite3 data/named-resource.db '.tables'

# Dump the resources table schema
sqlite3 data/named-resource.db '.schema resources'

# All rows, readable
sqlite3 -header -column data/named-resource.db \
  'SELECT uid, name, resource_type, resource_version FROM resources;'
```

The full resource spec is stored as JSON in the `spec` column. Use SQLite's
JSON functions to read individual fields:

```bash
# Extract spec fields from the JSON column
sqlite3 -header -column data/named-resource.db \
  "SELECT uid,
          json_extract(spec, '\$.name')   AS name,
          json_extract(spec, '\$.number') AS number
   FROM resources;"

# Find a row by a spec field
sqlite3 data/named-resource.db \
  "SELECT uid FROM resources WHERE json_extract(spec, '\$.name') = 'widget-a';"

# Count rows
sqlite3 data/named-resource.db 'SELECT COUNT(*) FROM resources;'
```

Note: in bash the `$` inside a JSON path (`$.name`) is escaped as `\$.name`
when the SQL is wrapped in double quotes, so the shell does not try to expand
it.

Write spec fields back with SQLite's JSON functions (they keep the `spec`
column valid JSON). The Named spec has `name`, `altName`, `number` and the
optional `somethingOrNothing`:

```bash
# Update a text field (match the row by another spec field)
sqlite3 data/named-resource.db \
  "UPDATE resources SET spec = json_set(spec, '\$.altName', 'widget-a-alt') WHERE json_extract(spec, '\$.name') = 'widget-a';"

# Update the numeric field (json('42') stays a number, not a quoted string)
sqlite3 data/named-resource.db \
  "UPDATE resources SET spec = json_set(spec, '\$.number', json('42')) WHERE json_extract(spec, '\$.name') = 'widget-a';"

# Add the optional field only if it is absent
sqlite3 data/named-resource.db \
  "UPDATE resources SET spec = json_insert(spec, '\$.somethingOrNothing', 'maybe') WHERE json_extract(spec, '\$.name') = 'widget-a';"

# Merge several fields at once (RFC 7396 patch)
sqlite3 data/named-resource.db \
  "UPDATE resources SET spec = json_patch(spec, json('{\"altName\":\"w-a\",\"number\":7}')) WHERE json_extract(spec, '\$.name') = 'widget-a';"
```

Direct writes bypass the API's validation and `resourceVersion` bumps — use
them for inspection/repair, not normal changes.

---

## Inspecting a Postgres database (psql)

If you run a Postgres-backed build instead (start it with, e.g.,
`--database-url "postgres://localhost/named-resource?sslmode=disable"`), the
same inspection is done with `psql`. The generic `resources` table stores the
spec as `jsonb`, so JSON access uses the `->` / `->>` operators instead of
`json_extract`.

```bash
PG="postgres://postgres:postgres@localhost:5432/named-resource?sslmode=disable"
```

No local `psql` install? Use the one bundled in the `postgres:16` image. If you
started Postgres as a named container (e.g. `docker run --name named-pg ...`),
exec into it:

```bash
docker exec -it named-pg psql -U postgres -d named-resource        # interactive
docker exec -i  named-pg psql -U postgres -d named-resource -c '\dt'  # one-off
```

Or run a throwaway client from the image (Docker Desktop resolves the host as
`host.docker.internal`; on Linux add `--add-host=host.docker.internal:host-gateway`
or use `--network host` with `localhost`):

```bash
docker run --rm -it postgres:16 \
  psql "postgres://postgres:postgres@host.docker.internal:5432/named-resource?sslmode=disable"
```

Interactive shell and meta-commands (the psql equivalents of sqlite's
dot-commands):

```bash
psql "$PG"
```

```
\dt            -- list tables            (sqlite: .tables)
\d resources  -- describe table+indexes  (sqlite: .schema / .indexes)
\di           -- list indexes
\x            -- toggle expanded output
\q            -- quit                    (sqlite: .quit)
```

One-off queries with `-c` (psql prints headers and aligns by default):

```bash
# List tables
psql "$PG" -c '\dt'

# Describe the resources table (columns, indexes, constraints)
psql "$PG" -c '\d resources'

# All rows
psql "$PG" -c 'SELECT uid, name, resource_type, resource_version FROM resources;'
```

Read individual spec fields from the `jsonb` column with `->` / `->>`:

```bash
# ->> returns text; cast when you need a number
psql "$PG" -c "SELECT uid,
                      spec->>'name'          AS name,
                      (spec->>'number')::int AS number
               FROM resources;"

# Find a row by a spec field
psql "$PG" -c "SELECT uid FROM resources WHERE spec->>'name' = 'widget-a';"

# Count rows
psql "$PG" -c 'SELECT COUNT(*) FROM resources;'

# Dump only the schema (like sqlite .schema) — uses pg_dump
pg_dump --schema-only --table=resources "$PG"
```

Write spec fields back with `jsonb_set` / the `||` merge operator (the value
must be valid JSON):

```bash
# Update a text field (a JSON string keeps its quotes: '\"...\"')
psql "$PG" -c "UPDATE resources SET spec = jsonb_set(spec, '{altName}', '\"widget-a-alt\"') WHERE spec->>'name' = 'widget-a';"

# Update the numeric field (written as JSON, not quoted text)
psql "$PG" -c "UPDATE resources SET spec = jsonb_set(spec, '{number}', '42') WHERE spec->>'name' = 'widget-a';"

# Merge several fields at once with ||
psql "$PG" -c "UPDATE resources SET spec = spec || '{\"altName\":\"w-a\",\"number\":7}'::jsonb WHERE spec->>'name' = 'widget-a';"

# Remove the optional field with the - operator
psql "$PG" -c "UPDATE resources SET spec = spec - 'somethingOrNothing' WHERE spec->>'name' = 'widget-a';"
```

Operator cheat sheet (SQLite → Postgres):

| Task              | SQLite                          | Postgres                     |
| ----------------- | ------------------------------- | ---------------------------- |
| List tables       | `.tables`                       | `\dt`                        |
| Table schema      | `.schema resources`             | `\d resources`               |
| List indexes      | `.indexes resources`            | `\d resources` or `\di`      |
| JSON field (text) | `json_extract(spec,'$.name')`   | `spec->>'name'`              |
| Nested JSON       | `json_extract(spec,'$.a.b')`    | `spec->'a'->>'b'`            |
| Headers/columns   | `-header -column`               | on by default; `\x` expanded |

Unlike SQLite (where `$` in a JSON path must be escaped), the Postgres `->>`
operator uses plain single-quoted keys, so no shell escaping is needed.
