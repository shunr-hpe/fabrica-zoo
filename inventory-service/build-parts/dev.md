# inventory-service — curl examples

These examples target the SMD-compatible resources generated under
`fabrica-zoo/inventory-service/inventory-service-sqlite/`. The service exposes
seven resource kinds, each backed by the generic Ent `resources` table
(SQLite). Authentication is disabled in the local build.

Build and run from inside the variant directory:

```bash
cd inventory-service-sqlite
go build -o bin/inventory-service ./cmd/server
rm -rf data; mkdir -p data
./bin/inventory-service serve --port 8080 \
  --database-url "file:./data/inventory.db?_fk=1" --auth-enabled=false
./bin/inventory-service serve --port 8080 \
  --database-url "file:./data/inventory.db?_fk=1"
```

Or just use the Makefile target, which builds and runs with the same flags:

```bash
make run
```

Set a base URL for convenience:

```bash
BASE=http://localhost:8080
```

---

## Resource overview

Each kind lives at its own collection. The spec fields follow SMD/Redfish
PascalCase naming. Resources are addressed by their generated `uid` (found in
`metadata.uid`), not by the SMD `ID`.

| Kind                | Collection            | Key spec fields                                                        |
| ------------------- | --------------------- | --------------------------------------------------------------------- |
| `Component`         | `/components`         | `ID`, `Type`, `State`, `Flag`, `Enabled`, `Role`, `SubRole`, `NID`, `Arch`, `Class` |
| `ComponentEndpoint` | `/componentendpoints` | `ID`, `Type`, `RedfishType`, `OdataID`, `RedfishEndpointID`, `Enabled` |
| `EthernetInterface` | `/ethernetinterfaces` | `ID`, `MACAddress`, `ComponentID`, `Type`, `IPAddresses`               |
| `Group`             | `/groups`             | `label`, `tags`, `members.ids`                                         |
| `Hardware`          | `/hardwares`          | `ID`, `Type`, `Ordinal`, `Status`, `HWInventoryByLocationType`         |
| `RedfishEndpoint`   | `/redfishendpoints`   | `ID`, `Type`, `FQDN`, `Hostname`, `Domain`, `Enabled`, `User`, `Password` |
| `ServiceEndpoint`   | `/serviceendpoints`   | `RedfishEndpointID`, `RedfishType`, `OdataID`, `UUID`                  |

Every resource shares the same envelope: a `metadata` block (with at least a
`name`) and a `spec` block. The server fills in `metadata.uid` and
`metadata.resourceVersion` on create.

---

## POST — create resources

### Component

```bash
curl -s -X POST "$BASE/components" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "x3000c0s1b0n0" },
    "spec": {
      "ID": "x3000c0s1b0n0",
      "Type": "Node",
      "State": "Ready",
      "Flag": "OK",
      "Enabled": true,
      "Role": "Compute",
      "NID": 1001,
      "Arch": "X86",
      "Class": "River"
    }
  }' | jq
```

The response includes the generated `metadata.uid`. Capture it for the
follow-up requests:

```bash
COMPONENT_UID=$(curl -s -X POST "$BASE/components" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "x3000c0s1b0n1" },
    "spec": {
      "ID": "x3000c0s1b0n1",
      "Type": "Node",
      "State": "On",
      "Enabled": true,
      "Role": "Compute",
      "NID": 1002,
      "Arch": "X86",
      "Class": "River"
    }
  }' | jq -r '.metadata.uid')

echo "Created COMPONENT_UID: $COMPONENT_UID"
```

### ComponentEndpoint

```bash
curl -s -X POST "$BASE/componentendpoints" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "x3000c0s1b0n0" },
    "spec": {
      "ID": "x3000c0s1b0n0",
      "Type": "Node",
      "RedfishType": "ComputerSystem",
      "RedfishSubtype": "Physical",
      "OdataID": "/redfish/v1/Systems/Node0",
      "RedfishEndpointID": "x3000c0s1b0",
      "Enabled": true
    }
  }' | jq
```

### EthernetInterface

The SMD `ID` is the MAC address with the separators removed.

```bash
curl -s -X POST "$BASE/ethernetinterfaces" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "a4bf0138e3e5" },
    "spec": {
      "ID": "a4bf0138e3e5",
      "MACAddress": "a4:bf:01:38:e3:e5",
      "ComponentID": "x3000c0s1b0n0",
      "Type": "Node",
      "IPAddresses": [
        { "IPAddress": "10.0.0.1", "Network": "NMN" }
      ]
    }
  }' | jq
```

### Group

```bash
curl -s -X POST "$BASE/groups" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "compute" },
    "spec": {
      "label": "compute",
      "tags": ["hpc"],
      "members": { "ids": ["x3000c0s1b0n0", "x3000c0s1b0n1"] }
    }
  }' | jq
```

### Hardware

```bash
curl -s -X POST "$BASE/hardwares" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "x3000c0s1b0n0" },
    "spec": {
      "ID": "x3000c0s1b0n0",
      "Type": "Node",
      "Ordinal": 0,
      "Status": "Populated",
      "HWInventoryByLocationType": "HWInvByLocNode"
    }
  }' | jq
```

### RedfishEndpoint

```bash
curl -s -X POST "$BASE/redfishendpoints" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "x3000c0s1b0" },
    "spec": {
      "ID": "x3000c0s1b0",
      "Type": "NodeBMC",
      "Hostname": "x3000c0s1b0",
      "Domain": "example.com",
      "FQDN": "x3000c0s1b0.example.com",
      "Enabled": true,
      "User": "root",
      "Password": "redfish-password",
      "RediscoverOnUpdate": true
    }
  }' | jq
```

### ServiceEndpoint

```bash
curl -s -X POST "$BASE/serviceendpoints" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "x3000c0s1b0-AccountService" },
    "spec": {
      "RedfishEndpointID": "x3000c0s1b0",
      "RedfishType": "AccountService",
      "OdataID": "/redfish/v1/AccountService",
      "UUID": "00000000-0000-0000-0000-000000000001"
    }
  }' | jq
```

---

## GET — list and fetch

List every resource in a collection:

```bash
curl -s "$BASE/components" | jq
```

Fetch a single resource by UID:

```bash
curl -s "$BASE/components/$COMPONENT_UID" | jq
```

The same pattern applies to every collection — swap `/components` for
`/componentendpoints`, `/ethernetinterfaces`, `/groups`, `/hardwares`,
`/redfishendpoints` or `/serviceendpoints`.

---

## PUT — update a resource

Replace the resource identified by `$COMPONENT_UID`. Include the full desired
spec.

```bash
curl -s -X PUT "$BASE/components/$COMPONENT_UID" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": { "name": "x3000c0s1b0n1" },
    "spec": {
      "ID": "x3000c0s1b0n1",
      "Type": "Node",
      "State": "Off",
      "Enabled": false,
      "Role": "Compute",
      "NID": 1002,
      "Arch": "X86",
      "Class": "River"
    }
  }' | jq
```

---

## DELETE — remove a resource

```bash
curl -s -X DELETE "$BASE/components/$COMPONENT_UID" | jq
```

If the response body is empty on delete, confirm removal with the status code:

```bash
curl -s -o /dev/null -w '%{http_code}\n' "$BASE/components/$COMPONENT_UID"
```

---

## Inspecting the SQLite database

The SQLite backend writes to the file named in `--database-url`
(`data/inventory.db` above). All seven kinds share one generic `resources`
table; rows are distinguished by the `kind` / `resource_type` column, and the
full spec is stored as JSON in the `spec` column.

No local `sqlite3` install? Run it from a container, mounting the project dir
so the `.db` file is visible. `alpine` is available everywhere (install sqlite
on the fly); `keinos/sqlite3` is a dedicated image with `sqlite3` as its
entrypoint:

```bash
# alpine (ubiquitous): install sqlite, then run
docker run --rm -it -v "$PWD:/work" -w /work alpine \
  sh -c 'apk add --no-cache sqlite >/dev/null && sqlite3 data/inventory.db'

# one-off query with alpine
docker run --rm -v "$PWD:/work" -w /work alpine \
  sh -c 'apk add --no-cache sqlite >/dev/null && sqlite3 data/inventory.db ".tables"'

# dedicated image (sqlite3 is the entrypoint — pass args directly)
docker run --rm -it -v "$PWD:/work" -w /work keinos/sqlite3 data/inventory.db
docker run --rm    -v "$PWD:/work" -w /work keinos/sqlite3 data/inventory.db '.tables'
```

To avoid reinstalling sqlite on every `alpine` run, build a reusable image once
(one-liner, Dockerfile piped over stdin) and reference it afterwards:

```bash
printf 'FROM alpine\nRUN apk add --no-cache sqlite\n' | docker build -t alpine-sqlite -

# then use it directly — no apk step needed
docker run --rm -it -v "$PWD:/work" -w /work alpine-sqlite sqlite3 data/inventory.db
docker run --rm    -v "$PWD:/work" -w /work alpine-sqlite sqlite3 data/inventory.db '.tables'
```

The examples below use a local `sqlite3`; prefix any of them with one of the
`docker run ...` wrappers above to run them from a container instead.

Open an interactive shell:

```bash
sqlite3 data/inventory.db
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
sqlite3 data/inventory.db '.tables'

# Dump the resources table schema
sqlite3 data/inventory.db '.schema resources'

# Count rows per kind
sqlite3 -header -column data/inventory.db \
  'SELECT kind, COUNT(*) AS n FROM resources GROUP BY kind;'

# All rows, readable (kind distinguishes the SMD resource types)
sqlite3 -header -column data/inventory.db \
  'SELECT uid, name, kind, resource_version FROM resources;'
```

The full resource spec is stored as JSON in the `spec` column. Use SQLite's
JSON functions to read individual fields. Filter by `kind` to scope a query to
one SMD resource type:

```bash
# Extract spec fields from Component rows
sqlite3 -header -column data/inventory.db \
  "SELECT uid,
          json_extract(spec, '\$.ID')    AS id,
          json_extract(spec, '\$.State') AS state,
          json_extract(spec, '\$.Role')  AS role
   FROM resources WHERE kind = 'Component';"

# Find a Component by its SMD ID
sqlite3 data/inventory.db \
  "SELECT uid FROM resources
   WHERE kind = 'Component' AND json_extract(spec, '\$.ID') = 'x3000c0s1b0n0';"

# Pull the MAC + component for every EthernetInterface
sqlite3 -header -column data/inventory.db \
  "SELECT json_extract(spec, '\$.MACAddress')  AS mac,
          json_extract(spec, '\$.ComponentID') AS component
   FROM resources WHERE kind = 'EthernetInterface';"

# Group members are a nested JSON array
sqlite3 data/inventory.db \
  "SELECT json_extract(spec, '\$.label') AS label,
          json_extract(spec, '\$.members.ids') AS member_ids
   FROM resources WHERE kind = 'Group';"
```

Note: in bash the `$` inside a JSON path (`$.ID`) is escaped as `\$.ID` when
the SQL is wrapped in double quotes, so the shell does not try to expand it.

Write spec fields back with SQLite's JSON functions (they keep the `spec`
column valid JSON):

```bash
# Update a text field (match the row by kind + SMD ID)
sqlite3 data/inventory.db \
  "UPDATE resources SET spec = json_set(spec, '\$.State', 'Off')
   WHERE kind = 'Component' AND json_extract(spec, '\$.ID') = 'x3000c0s1b0n0';"

# Set a boolean (json('false') keeps the JSON type instead of quoting it)
sqlite3 data/inventory.db \
  "UPDATE resources SET spec = json_set(spec, '\$.Enabled', json('false'))
   WHERE kind = 'Component' AND json_extract(spec, '\$.ID') = 'x3000c0s1b0n0';"

# Merge several fields at once (RFC 7396 patch)
sqlite3 data/inventory.db \
  "UPDATE resources SET spec = json_patch(spec, json('{\"Role\":\"Service\",\"NID\":2001}'))
   WHERE kind = 'Component' AND json_extract(spec, '\$.ID') = 'x3000c0s1b0n0';"
```

Direct writes bypass the API's validation and `resourceVersion` bumps — use
them for inspection/repair, not normal changes.
