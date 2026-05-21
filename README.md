# orders-api

REST API service built with Go, following Domain-Driven Design and Vertical Slice Architecture.

## Architecture

```
cmd/api/          → Composition root (pure DI, wires all layers)
config/           → Environment-based configuration
internal/
  domain/product/ → Aggregate root, Value Objects, Repository interface
  product/        → Vertical slices (one folder per use case)
    create/       → UseCase + Handler + Request + Response
    delete/       → UseCase + Handler
    find_by_id/   → UseCase + Handler + Response
    list/         → UseCase + Handler + Request + Response (paginated)
    update/       → UseCase + Handler + Request + Response (PATCH)
  infrastructure/
    postgres/product/ → Repository implementation + sqlc generated code
migrations/       → SQL migration files
server/           → Gin router, middleware, health check
```

**Key invariants:**
- All validations live in Value Objects (`domain/product/*.go`) — never in handlers
- Handlers only decode requests and encode responses
- Each use case injects only the dependencies it needs
- `internal/domain/product` has zero external imports

## Getting started

```bash
# 1. Install dependencies
go mod tidy

# 2. Start postgres
docker run -e POSTGRES_USER=dev -e POSTGRES_PASSWORD=dev -e POSTGRES_DB=dev \
  -p 5432:5432 postgres:16-alpine

# 3. Apply migrations
psql postgresql://dev:dev@localhost:5432/dev -f migrations/000001_create_products.up.sql

# 4. Run the server
DATABASE_URL=postgresql://dev:dev@localhost:5432/dev go run ./cmd/api
```

## API

| Method | Path                 | Description       |
|--------|----------------------|-------------------|
| POST   | /api/v1/products     | Create a product  |
| GET    | /api/v1/products     | List (paginated)  |
| GET    | /api/v1/products/:id | Get by ID         |
| PATCH  | /api/v1/products/:id | Partial update    |
| DELETE | /api/v1/products/:id | Delete            |
| GET    | /health              | Health check      |

### Pagination

```
GET /api/v1/products?page=1&page_size=20
```

### Partial update (PATCH)

Only provide the fields you want to change:

```json
{ "price": 49.99 }
```

## Testing

```bash
# Unit tests only
go test -short ./...

# Unit + integration tests (requires Docker)
go test ./...

# With coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Regenerate sqlc

```bash
sqlc generate
```

## Environment variables

| Variable     | Default | Description               |
|--------------|---------|---------------------------|
| DATABASE_URL | —       | PostgreSQL connection URL  |
| PORT         | 8080    | HTTP listen port           |
