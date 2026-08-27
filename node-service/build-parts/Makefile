# Copyright © 2025-2026 OpenCHAMI a Series of LF Projects, LLC
#
# SPDX-License-Identifier: MIT


# Service binary name
BINARY_NAME=named-resource-service
GO=go
GOFLAGS=-v

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build ldflags to embed version information
LDFLAGS=-ldflags "-X main.serviceVersion=$(VERSION) -X main.serviceCommit=$(COMMIT) -X main.serviceDate=$(DATE)"

GO_FILES=$(shell find . -name "*.go" -type f)
BINARY_DIR=bin
CONTAINER_CMD=docker
LOCAL_IMAGE_TAG ?= zoo-sqlite

# Default target
# Run all build targets
.PHONY: all
all: clean build

# Run all tests
.PHONY: all-tests
all-tests: test

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Build the binary
.PHONY: build
build: $(BINARY_DIR)/$(BINARY_NAME) ## Build the service binary

$(BINARY_DIR)/$(BINARY_NAME): $(GO_FILES)
	@echo "$(GO_FILES)"
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME) ./cmd/server/

.PHONY: build-client
build-client: $(BINARY_DIR)/$(BINARY_NAME)-client ## Build the client binary

$(BINARY_DIR)/$(BINARY_NAME)-client: $(GO_FILES)
	@echo "$(GO_FILES)"
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-client ./cmd/client/

.PHONY: goreleaser
goreleaser: GIT_STATE := $(shell if git diff-index --quiet HEAD --; then echo 'clean'; else echo 'dirty'; fi)
goreleaser: BUILD_HOST := $(shell hostname)
goreleaser: GO_VERSION := $(shell $(GO) version | awk '{print $$3}')
goreleaser: BUILD_USER := $(shell whoami)
goreleaser:
	@echo "GIT_STATE: $(GIT_STATE)"
	@echo "BUILD_HOST: $(BUILD_HOST)"
	@echo "GO_VERSION: $(GO_VERSION)"
	@echo "BUILD_USER: $(BUILD_USER)"
	@mkdir -p build
	@./scripts/generate-local-goreleaser .goreleaser.yaml build/.goreleaser-local.yaml
	goreleaser release --snapshot --clean --config build/.goreleaser-local.yaml
	@rm -f build/.goreleaser-local.yaml

# Clean build artifacts
.PHONY: clean
clean: ## Clean local build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf $(BINARY_DIR)
	rm -rf build
	rm -rf data
	rm -rf dist
	rm -f coverage.out coverage.html

.PHONY: clean-all
clean-all: clean ## Clean go and local build artifacts
	$(GO) clean --cache

.PHONY: run
run: build ## Build and run the service
	@mkdir -p data
	./bin/$(BINARY_NAME) serve --port 8080 --database-url file:data/named-resource.db?_fk=1 --auth-enabled=false

.PHONY: test
test: ## Run tests
	$(GO) test $(GOFLAGS) ./...

.PHONY: install
install: ## Install dependencies
	$(GO) mod download
	$(GO) mod verify

.PHONY: tidy
tidy: ## Tidy go.mod
	$(GO) mod tidy

.PHONY: fmt
fmt: ## Format code
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

# Build container image
.PHONY: image
image: build ## Build the container image
	@echo "Building container image..."
	# cp -r schemas bin
	# rm bin/schemas/*.go
	$(CONTAINER_CMD) build -f Dockerfile -t $(BINARY_NAME):$(LOCAL_IMAGE_TAG) bin

.PHONY: test
test: ## Run the unit tests
	$(GO) test -cover -v ./...

.DEFAULT_GOAL := help
