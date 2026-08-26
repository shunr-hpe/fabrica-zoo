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
UID=$(curl -s -X POST "$BASE/users" \
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

echo "Created UID: $UID"
```

---

## GET — list all users

```bash
curl -s "$BASE/users" | jq
```

## GET — fetch a single user by UID

```bash
curl -s "$BASE/users/$UID" | jq
```

Note: `password` is marked sensitive, so it is never returned in responses.

---

## PUT — update a user

Replace the user identified by `$UID`. Include the full desired spec.
`username` and `password` are immutable, so keep `username` unchanged.

```bash
curl -s -X PUT "$BASE/users/$UID" \
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
curl -s -X DELETE "$BASE/users/$UID" | jq
```

Confirm removal:

```bash
curl -s -o /dev/null -w '%{http_code}\n' "$BASE/users/$UID"
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
database directly, e.g. for SQLite:

```bash
sqlite3 data/user.db '.schema users' | grep -i spec_username
# -> spec_username text UNIQUE ...
```
