# inventory-service — build with `fabrica add resources`

These steps create a Fabrica service under
`fabrica-zoo/inventory-service/inventory-service/` from the hand-written spec
structs in [`structs/`](structs), using `fabrica add resources`.

Run them on Linux (bash). Paths assume you start in the repository root that
contains both `fabrica/` and `fabrica-zoo/`.

## Steps

```bash
# 0. Fabrica (and the generated service's codegen) require the Go 1.26.10
#    toolchain. This makes `go` auto-fetch it. Skip if your `go` is already
#    1.26.10 or newer.
export GOTOOLCHAIN=go1.26.10

# 1. Build the fabrica binary and put it on PATH.
cd fabrica
make clean build
export PATH=$(pwd)/bin:$PATH
cd ..

# 2. Initialize the project under fabrica-zoo/inventory-service/.
cd fabrica-zoo/inventory-service
fabrica init inventory-service \
  --module github.com/openchami/inventory-service \
  --group inventory-service.openchami.org \
  --storage-type ent \
  --db sqlite

# 3. Point the project at the local fabrica checkout.
cd inventory-service
printf '\nreplace github.com/openchami/fabrica => ../../../fabrica\n' >> go.mod

# 4. Add all resources from the struct files (one <Name>Spec per .go file).
#    ../structs is fabrica-zoo/inventory-service/structs.
fabrica add resources ../structs

# 5. Generate code, resolve dependencies, and build.
fabrica generate
go mod tidy
go build ./...
```

## Plain command list

```bash
export GOTOOLCHAIN=go1.26.10
cd fabrica
make clean build
export PATH=$(pwd)/bin:$PATH
cd ..
cd fabrica-zoo/inventory-service
fabrica init inventory-service --module github.com/openchami/inventory-service --group inventory-service.openchami.org --storage-type ent --db sqlite
cd inventory-service
printf '\nreplace github.com/openchami/fabrica => ../../../fabrica\n' >> go.mod
fabrica add resources ../structs
fabrica generate
go mod tidy
go build ./...
```

## Notes

- `fabrica add resources` must run from **inside** the project directory (it
  looks for `.fabrica.yaml` in the current directory). `../structs` resolves to
  `fabrica-zoo/inventory-service/structs`.
- The `replace ... => ../../../fabrica` line matches the existing
  `inventory-service-sqlite` variant and is only needed when building against the
  local fabrica checkout. Drop it to use a released fabrica.
- `GOTOOLCHAIN=go1.26.10` is required because the scaffolded `go.mod` declares
  `go 1.21`, so `GOTOOLCHAIN=auto` would not pull the 1.26.10 that the fabrica
  module requires.
