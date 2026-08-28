# user-service (Postgres) — dev notes

Postgres-focused companion to [`Dev.md`](./Dev.md). These notes cover starting a
Postgres instance and running a **Postgres-backed** `user-service` build end to
end. They target the two Postgres variants:

| Directory                  | fabrica base      | `spec.username` unique enforced? |
| -------------------------- | ----------------- | -------------------------------- |
| `user-service-postgres`    | promote-annotated | ✅ yes (dedicated `spec_username`) |
| `user-service-pr-postgres` | clean `pr-98`     | ❌ no (annotations dropped)        |

All commands assume the connection string is exported once:

```bash
export PG="postgres://postgres:postgres@localhost:5432/user_service?sslmode=disable"
```

---

## 1. Start Postgres

Fastest path — a throwaway container. `POSTGRES_DB` creates the `user_service`
database on first boot, so no manual `CREATE DATABASE` is needed:

```bash
docker run --rm -d --name user-pg -p 5432:5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=user_service \
  postgres:16
```

Wait until it is ready to accept connections, then confirm:

```bash
# poll readiness (exits 0 when accepting connections)
docker exec user-pg pg_isready -U postgres

# confirm the database exists
docker exec -i user-pg psql -U postgres -c '\l' | grep user_service
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
      POSTGRES_DB: user_service
    volumes:
      - user_pg_data:/var/lib/postgresql/data
volumes:
  user_pg_data:
```

Tear the throwaway container down when finished (this also drops all data):

```bash
docker rm -f user-pg
```

To create the database by hand against an existing server instead:

```bash
psql "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" \
  -c 'CREATE DATABASE user_service;'
```

---

## 2. Build and run the service

Build and run from inside the chosen variant directory:

```bash
cd user-service-pr-postgres        # or user-service-postgres
go build -o bin/user-service ./cmd/server
```

These services are generated with TokenSmith auth enabled. For local
development the simplest option is to disable auth on the command line:

```bash
./bin/user-service serve --port 8080 \
  --database-url "$PG" \
  --auth-enabled=false
```

To exercise auth instead of disabling it, leave `--auth-enabled` at its default
and point the server at a JWKS endpoint (AuthZ can stay off for CRUD testing):

```bash
export TOKENSMITH_JWKS_URL="http://localhost:3333/.well-known/jwks.json"
export TOKENSMITH_AUTHZ_MODE=off
./bin/user-service serve --port 8080 --database-url "$PG"
```

The database URL can also come from the environment instead of the flag; the
prefix is per-service:

```bash
export USER_SERVICE_PR_POSTGRES_DATABASE_URL="$PG"   # user-service-pr-postgres
# export USER_SERVICE_POSTGRES_DATABASE_URL="$PG"    # user-service-postgres
./bin/user-service serve --port 8080 --auth-enabled=false
```

On first start the service runs its Ent migrations and creates the tables in the
`user_service` database. Set a base URL for the curl examples:

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

Capture the generated `metadata.uid` for the follow-up requests:

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

`username` and `password` are immutable, so keep `username` unchanged:

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

The `username` spec field is declared `+fabrica:field:unique`. Posting two users
with the same `spec.username` must fail on the second attempt **on the
promote-annotated variant** (`user-service-postgres`); on the clean `pr-98`
variant (`user-service-pr-postgres`) both succeed because the annotation is not
emitted as a column constraint.

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

echo "Second POST status (expect 4xx/5xx on promote-annotated, 2xx on pr-98):"
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

---

## Inspecting the Postgres database (psql)

Use `psql` against the same `--database-url`. On the promote-annotated variant
(`user-service-postgres`) there is a dedicated `users` table with real `spec_*`
columns; on the clean `pr-98` variant (`user-service-pr-postgres`) the data
lives only in the generic `resources` table with a `jsonb` `spec` column.

No local `psql`? Use the one bundled in the running container, or a throwaway
client from the `postgres:16` image (Docker Desktop resolves the host as
`host.docker.internal`; on Linux add
`--add-host=host.docker.internal:host-gateway` or use `--network host` with
`localhost`):

```bash
# exec into the container started above
docker exec -it user-pg psql -U postgres -d user_service        # interactive
docker exec -i  user-pg psql -U postgres -d user_service -c '\dt'  # one-off

# throwaway client from the image
docker run --rm -it postgres:16 \
  psql "postgres://postgres:postgres@host.docker.internal:5432/user_service?sslmode=disable"
```

Interactive shell and meta-commands:

```bash
psql "$PG"
```

```
\dt            -- list tables
\d users       -- describe table + indexes/constraints
\di            -- list indexes
\x             -- toggle expanded (row-per-line) output
\q             -- quit
```

### Promote-annotated variant — dedicated `users` table

```bash
# columns, indexes and constraints (shows the promoted spec_* columns)
psql "$PG" -c '\d users'

# the unique constraint + composite idx_user_login index
psql "$PG" -c "SELECT indexname, indexdef FROM pg_indexes WHERE tablename='users';"

# query promoted columns directly (no JSON needed)
psql "$PG" -c 'SELECT uid, spec_username, spec_email, spec_role, spec_active FROM users;'

# non-promoted fields (e.g. fullName) stay in the spec_data jsonb blob
psql "$PG" -c "SELECT spec_username, spec_data->>'fullName' AS full_name FROM users;"

# duplicates on the unique column (returns nothing while the constraint holds)
psql "$PG" -c "SELECT spec_username, COUNT(*) AS n
               FROM users GROUP BY spec_username HAVING COUNT(*) > 1;"
```

### Clean `pr-98` variant — generic `resources` table

```bash
psql "$PG" -c '\d resources'

# everything is in the spec jsonb column; username is NOT a unique column here
psql "$PG" -c "SELECT uid,
                      spec->>'username' AS username,
                      spec->>'email'    AS email
               FROM resources;"

# duplicate usernames are allowed (no constraint) — returns rows here
psql "$PG" -c "SELECT spec->>'username' AS username, COUNT(*) AS n
               FROM resources GROUP BY 1 HAVING COUNT(*) > 1;"
```

### Writing values back

On the promote-annotated variant the `spec_*` fields are real columns; the
non-promoted `fullName` lives in the `spec_data` jsonb blob:

```bash
# promoted columns are ordinary columns
psql "$PG" -c "UPDATE users SET spec_role = 'admin', spec_active = true
               WHERE spec_username = 'alice';"

# fullName is not promoted — write it into spec_data with jsonb_set
psql "$PG" -c "UPDATE users SET spec_data = jsonb_set(spec_data, '{fullName}', '\"Alice Liddell\"')
               WHERE spec_username = 'alice';"
```

On the clean `pr-98` variant everything is in the generic `resources.spec` jsonb
column (`jsonb_set` for one field, `||` to merge, `-` to remove):

```bash
# update a text field (a JSON string keeps its quotes: '\"...\"')
psql "$PG" -c "UPDATE resources SET spec = jsonb_set(spec, '{role}', '\"admin\"')
               WHERE spec->>'username' = 'alice';"

# set a boolean (written as JSON, not quoted text)
psql "$PG" -c "UPDATE resources SET spec = jsonb_set(spec, '{active}', 'true')
               WHERE spec->>'username' = 'alice';"

# merge several fields at once
psql "$PG" -c "UPDATE resources SET spec = spec || '{\"role\":\"admin\",\"fullName\":\"Alice Liddell\"}'::jsonb
               WHERE spec->>'username' = 'alice';"
```

Caveats: direct writes bypass the API's validation and `resourceVersion` bumps.
On the dedicated table, writing a promoted column and its `spec_data` copy
separately can make them drift (the API keeps them in sync), and `password` is
stored bcrypt-hashed — never write `spec_password` as plaintext.

Unlike SQLite (where `$` in a JSON path must be escaped as `\$.username`), the
Postgres `->>` operator takes plain single-quoted keys, so no shell escaping is
needed inside a double-quoted `-c "..."`.
