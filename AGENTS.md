# AGENTS

This repository is a Go HTTP API backed by PostgreSQL.
Use this guide to keep changes consistent with existing patterns.

## Repository Overview
- Entry point: `cmd/http/main.go`
- HTTP server, routes, and middleware: `internal/httpserver`
- Feature packages: `internal/games`, `internal/accounts`
- Domain interfaces + errors: `internal/domain`
- Persistence helpers: `internal/storage/postgres`
- Shared utilities: `pkg/auth`, `pkg/httpjson`, `pkg/validator`
- SQL queries: `sql/queries`
- DB migrations: `sql/migrations`
- Generated sqlc code (read-only): `sql/sqlc_generated`

## Build / Run Commands
- Build all packages: `go build ./...`
- Run HTTP server: `go run ./cmd/http/main.go`
- Taskfile wrapper: `task serve` (runs server and pipes logs through `jq`)
- Install dev tooling: `task install-deps` (goose + sqlc)
- Dependency tidy: `go mod tidy`
- Dependency download: `go mod download`
- Go version: `1.25.3` (from `go.mod`)

## Lint / Format Commands
- Format code: `go fmt ./...`
- Vet code: `go vet ./...`
- No dedicated lint config found; add `golangci-lint` only if requested.

## Test Commands
- All tests: `go test ./...`
- Package tests: `go test ./internal/games`
- List tests in a package: `go test -list . ./internal/games`
- Single test (regex): `go test ./internal/games -run '^TestService_CreateGame$'`
- Single test (prefix): `go test ./internal/games -run '^TestHTTPAdapter_'`
- Disable cache during iteration: `go test ./internal/games -run TestService_CreateGame -count=1`
- DB-backed test example: `go test ./internal/games -run '^TestRepository_'`
- DB-backed test example: `go test ./internal/accounts -run '^TestRepository_'`
- DB-backed tests live in `internal/*/repository_test.go` (require Postgres).

## Database / Migrations
- Start local DB: `docker compose up -d`
- Adminer UI: `http://localhost:8080`
- Default dev DSN: `postgres://postgres:postgres@localhost:5432/nerd_backlog_dev?sslmode=disable`
- Run migrations: `task migrate`
- Create migration: `task add-migration -- <name>`
- Migrations use goose directives (`-- +goose Up/Down`, `StatementBegin/End`).
- Repository tests run migrations in `TestMain`.
- Reset DB volumes: `docker compose down -v` (destructive).

## SQLC
- Generate code: `task sqlc` (runs `sqlc generate`)
- Queries live in `sql/queries/*.sql`
- Use sqlc naming comments like `-- name: CreateGame :one`
- Use positional parameters (`$1`, `$2`, ...).
- Do not hand-edit `sql/sqlc_generated`

## Configuration
- Config is in `config/` with a `Development` env and default DSN.
- `config.NewHTTPConfig(config.Development)` is used in main and tests.
- If adding config, update `config/config.go` and `config/dev.go`.
- JWT manager is configured in `internal/httpserver/routes.go`.
- Current setup uses a static secret and TTLs; change with care.
- Avoid hardcoding secrets beyond local development defaults.

## Code Style Guidelines
- Keep files small and role-based (`service.go`, `repository.go`, `http_adapter.go`, `dto.go`).
- Prefer explicit interfaces in `internal/domain` and implementations per feature package.
- Use `context.Context` in service/repository signatures and pass `r.Context()` from handlers.
- Favor composition over large structs; keep adapters thin.

### Imports
- Group imports: stdlib, blank line, third-party, blank line, internal module.
- Use module path `github.com/kalogs-c/nerd-backlog/...` for internal imports.
- Use aliases when needed (e.g., `sqlc` for generated package).
- Let `go fmt` / `goimports` handle ordering.

### Naming
- Exported types/functions in PascalCase; unexported in lowerCamel.
- Constructors use `NewX` (e.g., `NewService`, `NewRepository`, `NewHTTPAdapter`).
- Errors are sentinel vars prefixed with `Err`.
- Keep receiver names short and consistent (`s`, `r`, `h`).
- Error strings are lower-case and descriptive (no trailing punctuation).

### Types and DTOs
- IDs are `uuid.UUID`; timestamps use `time.Time`.
- DTOs live in `dto.go` with JSON tags (snake_case).
- Mapping helpers live near DTOs (e.g., `MountAccountResponse`).
- Validation types implement `validator.Validator`.
- Use `validator.Problems` to return field-level issues.

### Error Handling
- Prefer explicit returns; avoid panics.
- Use `Must*` helpers only for truly fatal initialization (see `MustConnect`).
- Wrap errors with context using `fmt.Errorf("...: %w", err)`.
- Map `sql.ErrNoRows` to domain errors (e.g., `ErrGameNotFound`).
- Use `errors.Is` when comparing sentinel errors.

### HTTP Handlers
- Decode JSON with `httpjson.Decode` or `DecodeValid`.
- Return JSON with `httpjson.Encode`.
- Validation errors: `httpjson.EncodeValidationErrors`.
- Use `httpjson.NotifyError` / `NotifyHTTPError` for error responses + logging.
- Use explicit HTTP status codes: 400 for bad JSON, 404 for not found, 422 for invalid IDs.
- Prefer `httptest` for adapter tests.

### Logging
- Use `log/slog` with JSON handler in `cmd/http/main.go`.
- Middleware uses `WithLogging` to emit method/path/status/duration/remote.
- Prefer structured attributes (key/value pairs).
- Prefer `slog` over `log` except for fatal startup errors.

### Storage / Repositories
- Repositories depend on `*sqlc.Queries`.
- Keep repository methods thin; return domain models.
- Translate DB models into domain structs before returning.
- Use `pgxpool.Pool` via `internal/storage/postgres`.

### SQL and Migrations
- Queries should be minimal and named for sqlc (`-- name: Foo :one/:many/:exec`).
- Use `RETURNING *` when callers need full records.
- Migrations follow goose format with Up/Down blocks.
- Add indexes for lookup columns (see accounts email index).

### Testing
- Tests use `testify/require` for assertions.
- Use `testify/mock` for service/repository mocking.
- Integration tests depend on Postgres + migrations; keep data setup isolated.
- Use `uuid.New()` for test IDs and `mock.Anything` for dynamic args.
- Add `t.Helper()` in helpers you introduce.

### Generated Code
- `sql/sqlc_generated` is produced by sqlc; never edit by hand.
- If sqlc output changes, update repository usage accordingly.
- Regenerate after modifying `sql/queries`.

## Cursor / Copilot Rules
- No Cursor rules found in `.cursor/rules/` or `.cursorrules`.
- No Copilot instructions found in `.github/copilot-instructions.md`.

## Agent Workflow Tips
- Mirror existing patterns before introducing new abstractions.
- Prefer updating tests alongside behavior changes.
- Keep HTTP adapter and service responsibilities separate.
- Update DTOs, validators, and SQL together when adding endpoints.
- Run targeted tests first; expand to `go test ./...` when possible.
- Avoid touching generated code or migrations unless the change demands it.
