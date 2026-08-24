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
NAMED_UID=$(curl -s -X POST "$BASE/nameds" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "widget-b" },
    "spec": {
      "name": "widget-b",
      "altName": "Widget B",
      "number": 2
    }
  }' | jq -r '.metadata.uid')

echo "Created NAMED_UID: $NAMED_UID"
```

---

## GET — list all resources

```bash
curl -s "$BASE/nameds" | jq
```

## GET — fetch a single resource by UID

```bash
curl -s "$BASE/nameds/$NAMED_UID" | jq
```

---

## PUT — update a resource

Replace the resource identified by `$NAMED_UID`. Include the full desired spec.

```bash
curl -s -X PUT "$BASE/nameds/$NAMED_UID" \
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
curl -s -X DELETE "$BASE/nameds/$UID" | jq
```

If the response body is empty on delete, you can confirm removal with:

```bash
curl -s -o /dev/null -w '%{http_code}\n' "$BASE/nameds/$UID"
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
