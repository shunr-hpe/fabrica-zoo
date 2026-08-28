# named-resource-service (Postgres) — dev notes

Postgres-focused companion to [`Dev.md`](./Dev.md). These notes cover starting a
Postgres instance and running a **Postgres-backed** build of the service end to
end (start Postgres → run the service → CRUD → inspect with `psql`).

All commands assume the connection string is exported once:

```bash
export PG="postgres://postgres:postgres@localhost:5432/named_resource?sslmode=disable"
```

---

## 1. Start Postgres

Fastest path — a throwaway container. `POSTGRES_DB` creates the `named_resource`
database on first boot, so no manual `CREATE DATABASE` is needed:

```bash
docker run --rm -d --name named-pg -p 5432:5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=named_resource \
  postgres:16
```

Wait until it is ready to accept connections, then confirm:

```bash
# poll readiness (exits 0 when accepting connections)
docker exec named-pg pg_isready -U postgres

# confirm the database exists
docker exec -i named-pg psql -U postgres -c '\l' | grep named_resource
```

Prefer Compose? Drop this `docker-compose.yml` next to the service and
`docker compose up -d`:

```yaml
services:
  db:
    image: postgres:16
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: named_resource
    volumes:
      - named_pg_data:/var/lib/postgresql/data
volumes:
  named_pg_data:
```

Tear the throwaway container down when finished (this also drops all data):

```bash
docker rm -f named-pg
```

To create the database by hand against an existing server instead:

```bash
psql "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" \
  -c 'CREATE DATABASE named_resource;'
```

---

## 2. Build and run the service

Build and run from inside the Postgres variant directory:

```bash
cd node-service-postgres
go build -o bin/named-resource-service ./cmd/server
```

This service is generated with TokenSmith auth enabled. For local development
the simplest option is to disable auth on the command line:

```bash
./bin/named-resource-service serve --port 8080 \
  --database-url "$PG" \
  --auth-enabled=false
```

To exercise auth instead of disabling it, leave `--auth-enabled` at its default
and point the server at a JWKS endpoint (AuthZ can stay off for CRUD testing):

```bash
export TOKENSMITH_JWKS_URL="http://localhost:3333/.well-known/jwks.json"
export TOKENSMITH_AUTHZ_MODE=off
./bin/named-resource-service serve --port 8080 --database-url "$PG"
```

The database URL can also come from the environment instead of the flag (the
prefix is the upper-snake-case service name, e.g.
`NODE_SERVICE_POSTGRES_DATABASE_URL`).

On first start the service runs its Ent migrations and creates the tables in the
`named_resource` database. Set a base URL for the curl examples:

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

Capture the generated `metadata.uid` for the follow-up requests:

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

Confirm removal:

```bash
curl -s -o /dev/null -w '%{http_code}\n' "$BASE/nameds/$RESOURCE_UID"
```

---

## Name uniqueness enforcement test

The `name` spec field is declared unique. Posting two resources with the same
`spec.name` must fail on the second attempt.

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

## Inspecting the Postgres database (psql)

The generic backend stores the spec as a `jsonb` column in the `resources`
table, so JSON access uses the `->` / `->>` operators instead of SQLite's
`json_extract`.

No local `psql`? Use the one bundled in the running container, or a throwaway
client from the `postgres:16` image (Docker Desktop resolves the host as
`host.docker.internal`; on Linux add
`--add-host=host.docker.internal:host-gateway` or use `--network host` with
`localhost`):

```bash
# exec into the container started above
docker exec -it named-pg psql -U postgres -d named_resource        # interactive
docker exec -i  named-pg psql -U postgres -d named_resource -c '\dt'  # one-off

# throwaway client from the image
docker run --rm -it postgres:16 \
  psql "postgres://postgres:postgres@host.docker.internal:5432/named_resource?sslmode=disable"
```

Interactive shell and meta-commands (the psql equivalents of sqlite's
dot-commands):

```bash
psql "$PG"
```

```
\dt            -- list tables            (sqlite: .tables)
\d resources   -- describe table+indexes (sqlite: .schema / .indexes)
\di            -- list indexes
\x             -- toggle expanded output
\q             -- quit                    (sqlite: .quit)
```

One-off queries with `-c` (psql prints headers and aligns by default):

```bash
# list tables
psql "$PG" -c '\dt'

# describe the resources table (columns, indexes, constraints)
psql "$PG" -c '\d resources'

# all rows
psql "$PG" -c 'SELECT uid, name, resource_type, resource_version FROM resources;'
```

Read individual spec fields from the `jsonb` column with `->` / `->>`:

```bash
# ->> returns text; cast when you need a number
psql "$PG" -c "SELECT uid,
                      spec->>'name'          AS name,
                      (spec->>'number')::int AS number
               FROM resources;"

# find a row by a spec field
psql "$PG" -c "SELECT uid FROM resources WHERE spec->>'name' = 'widget-a';"

# count rows
psql "$PG" -c 'SELECT COUNT(*) FROM resources;'

# dump only the schema (like sqlite .schema) — uses pg_dump
pg_dump --schema-only --table=resources "$PG"
```

Write spec fields back with `jsonb_set` / the `||` merge operator (the value
must be valid JSON):

```bash
# update a text field (a JSON string keeps its quotes: '\"...\"')
psql "$PG" -c "UPDATE resources SET spec = jsonb_set(spec, '{altName}', '\"widget-a-alt\"') WHERE spec->>'name' = 'widget-a';"

# update the numeric field (written as JSON, not quoted text)
psql "$PG" -c "UPDATE resources SET spec = jsonb_set(spec, '{number}', '42') WHERE spec->>'name' = 'widget-a';"

# merge several fields at once with ||
psql "$PG" -c "UPDATE resources SET spec = spec || '{\"altName\":\"w-a\",\"number\":7}'::jsonb WHERE spec->>'name' = 'widget-a';"

# remove the optional field with the - operator
psql "$PG" -c "UPDATE resources SET spec = spec - 'somethingOrNothing' WHERE spec->>'name' = 'widget-a';"
```

Direct writes bypass the API's validation and `resourceVersion` bumps — use them
for inspection/repair, not normal changes.

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
