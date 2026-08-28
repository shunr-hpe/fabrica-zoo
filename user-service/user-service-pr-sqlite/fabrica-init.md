# user-service-pr-sqlite generation

This project was generated with `fabrica-init`, driving the local `pr-98` branch
of `fabrica`, with TokenSmith authentication enabled and SQLite storage.

## Prerequisites

The local `fabrica` `pr-98` branch requires Go 1.26.6 (matching the system Go),
so the toolchain is pinned for the generate/tidy steps:

```powershell
$env:GOTOOLCHAIN = "go1.26.6"
```

The `fabrica` CLI was built from the local `pr-98` branch:

```powershell
# in fabrica/
go build -o bin/fabrica.exe ./cmd/fabrica
```

The `fabrica-init` helper was built from source:

```powershell
# in fabrica-init/
go build -o bin/fabrica-init.exe .
```

## Command

Run from the repository root (`p/`):

```powershell
$env:GOTOOLCHAIN = "go1.26.6"

.\fabrica-init\bin\fabrica-init.exe `
  --module github.com/openchami/user-service-pr-sqlite `
  --group storage-annotations.openchami.org `
  --db sqlite `
  --auth `
  --name user-service-pr-sqlite `
  --dir fabrica-zoo/user-service `
  --local-fabrica ../../../fabrica `
  --fabrica ~/work/p/fabrica/bin/fabrica.exe `
  --resources fabrica-zoo/user-service/structs
```

- `--auth` enables authentication with TokenSmith (generates the AuthZ
  classifier and starter `authz/` policy files).
- `--db sqlite` selects the SQLite Ent driver.
- `--local-fabrica ../../../fabrica` adds a `replace` directive in the generated
  `go.mod` pointing at the local `fabrica` checkout.
- `--resources fabrica-zoo/user-service/structs` derives the `User` resource
  from the `UserSpec` struct in `structs/user.go`.

## Output

```text
+ ~/work/p/fabrica/bin/fabrica.exe init user-service-pr-sqlite --module github.com/openchami/user-service-pr-sqlite --group storage-annotations.openchami.org --storage-version v1 --storage-type ent --db sqlite --auth  (in fabrica-zoo/user-service)
🚀 Creating user-service-pr-sqlite project...
  ├─ Created .fabrica.yaml
  ├─ Created apis.yaml (group storage-annotations.openchami.org, storage v1)

✅ Project initialized successfully!

Next steps:
  1. Add resources with 'fabrica add resource <name>'
  2. Define types in apis/storage-annotations.openchami.org/<version>/*_types.go
  3. Run 'fabrica generate' to generate code
  4. Run 'go mod tidy' to update dependencies
  5. Start development with 'go run ./cmd/server/'

+ append to go.mod:
replace github.com/openchami/fabrica => ../../../fabrica
+ ~/work/p/fabrica/bin/fabrica.exe add resource User  (in fabrica-zoo\user-service\user-service-pr-sqlite)
No version specified, using storage hub version: v1
📦 Adding resource User to storage-annotations.openchami.org/v1...
  ✓ Added User to apis.yaml

✅ Resource added successfully!

Next steps:
  1. Edit apis\storage-annotations.openchami.org\v1\user_types.go to customize your resource
  2. Add to other versions with 'fabrica add version <new-version>'
  3. Run 'fabrica generate' to create handlers

+ merged UserSpec into apis\storage-annotations.openchami.org\v1\user_types.go
+ ~/go/bin/goimports.exe -w apis\storage-annotations.openchami.org\v1\user_types.go  (in fabrica-zoo\user-service\user-service-pr-sqlite)
+ ~/work/p/fabrica/bin/fabrica.exe generate  (in fabrica-zoo\user-service\user-service-pr-sqlite)
🔧 Generating code...
📦 Found 1 resource(s): User

📝 Registration file not found, creating it...
🔍 Discovering resources...
📦 Found 1 resource(s): User

✅ Generated pkg\resources\register_generated.go

go: found github.com/openchami/fabrica/pkg/codegen in github.com/openchami/fabrica v0.0.0-00010101000000-000000000000
go: found github.com/openchami/fabrica/pkg/fabrica in github.com/openchami/fabrica v0.0.0-00010101000000-000000000000
🛠️  Generating handlers...
  ✓ Generated cmd\server\user_handlers_generated.go
⚙️  Generating middleware...
  ✓ Generated internal\middleware\validation_middleware_generated.go
  ✓ Generated internal\middleware\conditional_middleware_generated.go
  ✓ Generated internal\middleware\versioning_middleware_generated.go
📊 Generating metrics stub (metrics disabled)...
  ✓ Generated cmd\server\metrics_helpers_generated.go
🗄️  Generating Ent schemas...
  ✓ Generated internal\storage\ent\schema\resource.go
  ✓ Generated internal\storage\ent\schema\label.go
  ✓ Generated internal\storage\ent\schema\annotation.go
🔗 Generating Ent adapters...
  ✓ Generated internal\storage\ent_adapter.go
  ✓ Generated internal\storage\generate.go
🧰 Generating Ent helpers (queries, transactions)...
  ✓ Generated internal\storage\ent_queries_generated.go
  ✓ Generated internal\storage\ent_transactions_generated.go
📁 Generating storage layer (ent)...
  ✓ Generated internal\storage\storage_generated.go
📋 Generating OpenAPI specification...
  ✓ Generated cmd\server\openapi_generated.go
  ✓ Generated cmd\server\openapi_extensions.go
  ✓ Generated cmd\server\version_generated.go
🛣️  Generating routes...
  ✓ Generated cmd\server\routes_generated.go
📊 Generating models...
  ✓ Generated cmd\server\models_generated.go
🔐 Generating AuthZ classifier...
  ✓ Generated cmd\server\authz_classifier_generated.go
  ✓ Generated cmd\server\authz_classifier.go
📜 Generating starter AuthZ policy files...
  ✓ Generated authz\model.conf
  ✓ Generated authz\policy.csv
  ✓ Generated authz\grouping.csv
📤 Generating export command...
  ✓ Generated cmd\server\export.go
📥 Generating import command...
  ✓ Generated cmd\server\import.go
🔄 Generating API version registry...
  ✓ Generated pkg\apiversion\registry_generated.go
📦 Generating client code...
🔌 Generating client library...
  ✓ Generated pkg\client\client_generated.go
📊 Generating client models...
  ✓ Generated pkg\client\models_generated.go
⚡ Generating CLI client...
  ✓ Generated cmd\client\main.go
  ✓ Generated cmd\client\version_generated.go
🔄 Generating Ent client code...
✅ Ent client code generated
  └─ Done!

✅ Code generation complete!

Next steps:
  go mod tidy                     # Update dependencies
  go run ./cmd/server       # Start the server
+ go mod tidy  (in fabrica-zoo\user-service\user-service-pr-sqlite)
go: finding module for package github.com/getkin/kin-openapi/openapi3
go: finding module for package github.com/getkin/kin-openapi/openapi3gen
go: downloading github.com/getkin/kin-openapi v0.149.0
go: found github.com/getkin/kin-openapi/openapi3 in github.com/getkin/kin-openapi v0.149.0
go: found github.com/getkin/kin-openapi/openapi3gen in github.com/getkin/kin-openapi v0.149.0

EXITCODE=0
```

## Verification

The generated project compiles cleanly:

```powershell
# in fabrica-zoo/user-service/user-service-pr-sqlite/
$env:GOTOOLCHAIN = "go1.26.6"
go build ./...   # exit 0
```
