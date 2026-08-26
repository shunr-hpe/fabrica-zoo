# user-service variants — curl examples

These examples target the `User` resource generated under
`fabrica-zoo/storage-annotations/`. There are four variants:

| Directory                  | DB       | fabrica base        | `spec.username` unique enforced? |
| -------------------------- | -------- | ------------------- | -------------------------------- |
| `user-service`             | SQLite   | promote-annotated   | ✅ yes (dedicated `spec_username`) |
| `user-service-postgres`    | Postgres | promote-annotated   | ✅ yes (dedicated `spec_username`) |
| `user-service-pr-sqlite`   | SQLite   | clean `pr-106`      | ❌ no (annotations dropped)        |
| `user-service-pr-postgres` | Postgres | clean `pr-106`      | ❌ no (annotations dropped)        |

The uniqueness test at the bottom is the point of the comparison: the
promote-annotated variants reject a duplicate `spec.username` at the database
layer, while the clean `pr-106` variants accept it because the
`+fabrica:field:unique` annotation is not emitted as a column constraint.

---

## Running a variant

Build and run from inside the variant directory. Pick the matching
`--database-url` for the backend.

SQLite variant (`user-service` or `user-service-pr-sqlite`):

```bash
cd user-service            # or user-service-pr-sqlite
go build -o bin/user-service ./cmd/server
rm -rf data; mkdir -p data
./bin/user-service serve --port 8080 --database-url "file:./data/user.db?_fk=1"
```

Postgres variant (`user-service-postgres` or `user-service-pr-postgres`):

```bash
cd user-service-postgres   # or user-service-pr-postgres
go build -o bin/user-service ./cmd/server

# Start a throwaway Postgres (or point at your own):
docker run --rm -d --name user-pg -p 5432:5432 \
  -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=user_service postgres:16

./bin/user-service serve --port 8080 \
  --database-url "postgres://postgres:postgres@localhost:5432/user_service?sslmode=disable"
```

Set a base URL for convenience:

```bash
BASE=http://localhost:8080
```

The `User` resource lives at the `/users` collection. Its spec fields are:

| Field      | Type   | Notes                                                    |
| ---------- | ------ | -------------------------------------------------------- |
| `username` | string | **unique**, indexed, immutable, `min=3,max=32,alphanum`  |
| `email`    | string | **unique**, indexed, `email` format                      |
| `password` | string | bcrypt-hashed, immutable, sensitive, `min=8`             |
| `fullName` | string | required, `min=1,max=128` (stays in `spec_data` JSON)    |
| `role`     | string | default `user`, indexed, `oneof=admin user readonly`     |
| `active`   | bool   | default `true`                                           |

Resources are addressed by their generated `uid` (found in `metadata.uid`),
not by `username`.

---

## POST — create a user

```bash
curl -s -X POST "$BASE/users" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "alice" },
    "spec": {
      "username": "alice",
      "email": "alice@example.com",
      "password": "s3cretpass",
      "fullName": "Alice Example",
      "role": "admin",
      "active": true
    }
  }' | jq
```

The response includes the generated `metadata.uid`. Capture it for the
follow-up requests:

```bash
RESOURCE_UID=$(curl -s -X POST "$BASE/users" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "bob" },
    "spec": {
      "username": "bob",
      "email": "bob@example.com",
      "password": "s3cretpass",
      "fullName": "Bob Example",
      "role": "user"
    }
  }' | jq -r '.metadata.uid')

echo "Created RESOURCE_UID: $RESOURCE_UID"
```

---

## GET — list all users

```bash
curl -s "$BASE/users" | jq
```

## GET — fetch a single user by UID

```bash
curl -s "$BASE/users/$RESOURCE_UID" | jq
```

Note: `password` is marked sensitive, so it is never returned in responses.

---

## PUT — update a user

Replace the user identified by `$RESOURCE_UID`. Include the full desired spec.
`username` and `password` are immutable, so keep `username` unchanged.

```bash
curl -s -X PUT "$BASE/users/$RESOURCE_UID" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "bob" },
    "spec": {
      "username": "bob",
      "email": "bob.updated@example.com",
      "password": "s3cretpass",
      "fullName": "Bob Example (updated)",
      "role": "readonly",
      "active": false
    }
  }' | jq
```

---

## DELETE — remove a user

```bash
curl -s -X DELETE "$BASE/users/$RESOURCE_UID" | jq
```

Confirm removal:

```bash
curl -s -o /dev/null -w '%{http_code}\n' "$BASE/users/$RESOURCE_UID"
```

---

## Username uniqueness enforcement test

The `username` spec field is declared `+fabrica:field:unique`. Posting two
users with the same `spec.username` must fail on the second attempt **on the
promote-annotated variants**.

```bash
# First create with username "carol" — expected to succeed.
echo "First POST (should succeed):"
curl -s -X POST "$BASE/users" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "carol-1" },
    "spec": {
      "username": "carol",
      "email": "carol1@example.com",
      "password": "s3cretpass",
      "fullName": "Carol One",
      "role": "user"
    }
  }' | jq

# Second create with the same spec.username "carol" but a different email —
# expected to FAIL on the promote-annotated variants.
echo "Second POST (should fail on uniqueness):"
curl -s -X POST "$BASE/users" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "carol-2" },
    "spec": {
      "username": "carol",
      "email": "carol2@example.com",
      "password": "s3cretpass",
      "fullName": "Carol Two",
      "role": "user"
    }
  }' | jq
```

To make the outcome explicit, capture the HTTP status codes:

```bash
echo "First POST status:"
curl -s -o /dev/null -w '%{http_code}\n' -X POST "$BASE/users" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "dave-1" },
    "spec": {
      "username": "dave",
      "email": "dave1@example.com",
      "password": "s3cretpass",
      "fullName": "Dave One",
      "role": "user"
    }
  }'

echo "Second POST status:"
curl -s -o /dev/null -w '%{http_code}\n' -X POST "$BASE/users" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "dave-2" },
    "spec": {
      "username": "dave",
      "email": "dave2@example.com",
      "password": "s3cretpass",
      "fullName": "Dave Two",
      "role": "user"
    }
  }'
```

Expected results:

- **`user-service` / `user-service-postgres`** (promote-annotated): the first
  POST returns `2xx`; the second returns a `4xx/5xx` error because the unique
  constraint on the `spec_username` column rejects the duplicate.
- **`user-service-pr-sqlite` / `user-service-pr-postgres`** (clean `pr-106`):
  both POSTs return `2xx` — the `unique` annotation was dropped, so duplicate
  usernames are stored in the generic JSON `spec` column with no constraint.

You can prove the difference on a promote-annotated variant by inspecting the
database directly. See the next section.

---

## Inspecting the SQLite database

Applies to the SQLite variants (`user-service`, `user-service-pr-sqlite`),
which write to `data/user.db`. (Postgres variants use `psql` — see the end of
this section.)

Open an interactive shell and orient yourself:

```bash
sqlite3 data/user.db
```

```
.tables       -- list tables
.schema       -- full schema for every table
.headers on   -- show column names
.mode column  -- align output
.quit         -- exit
```

On `user-service` (promote-annotated) there is a dedicated `users` table with
real columns; on `user-service-pr-sqlite` the data lives only in the generic
`resources` table.

### Promote-annotated variant — dedicated `users` table

```bash
# Schema shows the promoted spec_* columns and their constraints
sqlite3 data/user.db '.schema users'

# The unique constraint on spec_username
sqlite3 data/user.db '.schema users' | grep -i spec_username
# -> spec_username text UNIQUE ...

# List the indexes on the table (includes idx_user_login composite unique)
sqlite3 data/user.db '.indexes users'
sqlite3 -header -column data/user.db \
  "SELECT name, sql FROM sqlite_master WHERE type='index' AND tbl_name='users';"

# Query promoted columns directly (no JSON needed)
sqlite3 -header -column data/user.db \
  'SELECT uid, spec_username, spec_email, spec_role, spec_active FROM users;'

# Non-promoted fields (e.g. fullName) stay in the spec_data JSON blob
sqlite3 -header -column data/user.db \
  "SELECT spec_username, json_extract(spec_data, '\$.fullName') AS full_name FROM users;"
```

### Base variant (`user-service-pr-sqlite`) — generic `resources` table

```bash
sqlite3 data/user.db '.schema resources'

# Everything is in the spec JSON column; username is NOT a unique column here
sqlite3 -header -column data/user.db \
  "SELECT uid,
          json_extract(spec, '\$.username') AS username,
          json_extract(spec, '\$.email')    AS email
   FROM resources;"

# Duplicate usernames are allowed (no constraint) — this returns rows here,
# but nothing on the promote-annotated variant
sqlite3 data/user.db \
  "SELECT json_extract(spec, '\$.username') AS username, COUNT(*) AS n
   FROM resources GROUP BY username HAVING n > 1;"
```

Note: in bash the `$` inside a JSON path (`$.username`) is escaped as
`\$.username` when the SQL is wrapped in double quotes.

### Inspecting a Postgres database (psql)

For the Postgres variants (`user-service-postgres`, `user-service-pr-postgres`)
use `psql` against the same `--database-url`. Promoted `spec_*` fields are real
columns; JSON fields live in a `jsonb` column and use the `->` / `->>`
operators instead of `json_extract`.

```bash
PG="postgres://postgres:postgres@localhost:5432/user_service?sslmode=disable"
```

Meta-commands (psql equivalents of the sqlite dot-commands):

```
\dt        -- list tables                         (sqlite: .tables)
\d users   -- describe table + indexes/constraints (sqlite: .schema / .indexes)
\di        -- list indexes
\x         -- expanded row output
\q         -- quit                                 (sqlite: .quit)
```

Promote-annotated variant (`user-service-postgres`) — dedicated `users` table:

```bash
# Columns, constraints and indexes: shows the UNIQUE on spec_username and the
# idx_user_login composite unique index
psql "$PG" -c '\d users'

# Just the indexes/constraints
psql "$PG" -c "SELECT indexname, indexdef FROM pg_indexes WHERE tablename='users';"

# Query promoted columns directly (no JSON needed)
psql "$PG" -c 'SELECT uid, spec_username, spec_email, spec_role, spec_active FROM users;'

# Non-promoted fields (e.g. fullName) live in the spec_data jsonb column
psql "$PG" -c "SELECT spec_username, spec_data->>'fullName' AS full_name FROM users;"
```

Base variant (`user-service-pr-postgres`) — generic `resources` table:

```bash
psql "$PG" -c '\d resources'

# Everything is in the spec jsonb column; username is NOT a unique column here
psql "$PG" -c "SELECT uid, spec->>'username' AS username, spec->>'email' AS email FROM resources;"

# Duplicate usernames are allowed (no constraint) — returns rows here, but
# nothing on the promote-annotated variant
psql "$PG" -c "SELECT spec->>'username' AS username, COUNT(*) AS n
               FROM resources GROUP BY 1 HAVING COUNT(*) > 1;"
```

Operator cheat sheet (SQLite → Postgres):

| Task              | SQLite                                    | Postgres                     |
| ----------------- | ----------------------------------------- | ---------------------------- |
| List tables       | `.tables`                                 | `\dt`                        |
| Table schema      | `.schema users`                           | `\d users`                   |
| List indexes      | `.indexes users`                          | `\d users` or `\di`          |
| JSON field (text) | `json_extract(spec_data,'$.fullName')`    | `spec_data->>'fullName'`     |
| Nested JSON       | `json_extract(col,'$.a.b')`               | `col->'a'->>'b'`             |
| Headers/columns   | `-header -column`                         | on by default; `\x` expanded |


