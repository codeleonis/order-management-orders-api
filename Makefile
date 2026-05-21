## ─── Variables ────────────────────────────────────────────────────────────────
BINARY        := bin/server
DATABASE_URL  ?= postgresql://dev:dev@localhost:5432/dev?sslmode=disable
IMAGE_NAME    ?= your-service
COVERAGE_FILE := coverage.out

## ─── Build ────────────────────────────────────────────────────────────────────
.PHONY: build
build:                          ## Compile the binary to ./bin/server
	@mkdir -p bin
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY) ./cmd/api

.PHONY: run
run:                            ## Run the server locally (requires DATABASE_URL)
	DATABASE_URL=$(DATABASE_URL) go run ./cmd/api

.PHONY: tidy
tidy:                           ## Tidy and verify go modules
	go mod tidy
	go mod verify

## ─── Tests ────────────────────────────────────────────────────────────────────
.PHONY: test
test:                           ## Run unit tests only (no Docker required)
	go test -short -v -race ./...

.PHONY: test/all
test/all:                       ## Run unit + integration tests (requires Docker)
	go test -v -race ./...

.PHONY: test/coverage
test/coverage:                  ## Run unit tests and open HTML coverage report
	go test -short -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@COVERAGE=$$(go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print substr($$3, 1, length($$3)-1)}'); \
	echo "Total coverage: $${COVERAGE}%"; \
	awk "BEGIN { if ($${COVERAGE} < 80) { print \"Coverage $${COVERAGE}% is below 80%\"; exit 1 } }"
	go tool cover -html=$(COVERAGE_FILE)

## ─── Code quality ─────────────────────────────────────────────────────────────
.PHONY: lint
lint:                           ## Run golangci-lint
	golangci-lint run ./...

.PHONY: vuln
vuln:                           ## Run govulncheck
	govulncheck ./...

.PHONY: check
check: lint vuln test           ## Run lint + vuln + unit tests

## ─── Database ─────────────────────────────────────────────────────────────────
.PHONY: migrate/up
migrate/up:                     ## Apply all migrations
	psql "$(DATABASE_URL)" -f migrations/000001_create_products.up.sql

.PHONY: migrate/down
migrate/down:                   ## Rollback all migrations
	psql "$(DATABASE_URL)" -f migrations/000001_create_products.down.sql

## ─── Code generation ──────────────────────────────────────────────────────────
.PHONY: sqlc
sqlc:                           ## Regenerate sqlc query code
	sqlc generate

## ─── Docker ───────────────────────────────────────────────────────────────────
.PHONY: docker/build
docker/build:                   ## Build the Docker image
	docker build -t $(IMAGE_NAME) .

.PHONY: docker/up
docker/up:                      ## Start app + postgres via docker compose
	docker compose up --build

.PHONY: docker/down
docker/down:                    ## Stop and remove containers
	docker compose down

.PHONY: docker/logs
docker/logs:                    ## Tail app logs
	docker compose logs -f app

## ─── Help ─────────────────────────────────────────────────────────────────────
.PHONY: help
help:                           ## Show this help
	@grep -E '^[a-zA-Z_/%-]+:.*##' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*##"}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
