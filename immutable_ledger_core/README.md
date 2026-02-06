# Immutable Ledger Core

Lightweight Go service demonstrating an immutable ledger-style transaction store.

Repository layout

- `cmd/main.go` — application entrypoint.
- `app/transactions/` — transaction domain logic and HTTP handlers:
  - `transaction_database.go` — in-memory store (likely `TranasctionDatabase`).
  - `transaction_model.go` — domain model(s) (Transaction, DTOs).
  - `transaction_service.go` — business logic (create/validate transactions).
  - `transaction_routes.go` — Gin HTTP handlers (create, list, validate).
  - `transaction_error.go` — typed errors (e.g., `TransactionError`).
- `app/utils/` — helper utilities:
  - `id_generator.go` — ID generation helper(s).
  - `hash_generator.go` — hashing helper(s).
- `go.mod`, `Makefile` — build metadata.

Project purpose

This repo is a simple example of a small service that accepts transactions, stores them in an (in-memory) append-only ledger, and can validate the ledger integrity. It’s useful for learning concepts around immutability, hashing, and service wiring in Go.

Implemented features

- HTTP API (Gin) with handlers for:
  - Creating a transaction (accepts JSON DTO)
  - Listing all transactions
  - Validating the ledger
- In-memory transaction store (simple append-only store)
- Pluggable ID and hash generators (handlers accept generator functions)
- Typed domain error (`TransactionError`) used to communicate business errors

What’s missing / recommended improvements

Tests
- Unit tests: None present. Add unit tests for:
  - `CreateTransaction` happy path and edge cases (invalid input, duplicate id, hash mismatch)
  - `ValidateTransactions` with valid and tampered ledger
  - Error types and conversion (e.g., verify HTTP handlers map typed errors to proper codes)
- Integration tests: HTTP tests using httptest and Gin router to confirm endpoints and JSON shapes.

Error handling and response flow
- Some handlers currently may try to write multiple responses on error paths (see `transaction_routes.go` where a type-assert branch writes a JSON but the function continues and writes a 500 too). Ensure each branch returns after writing a response.
- Consider using `errors.Is` / `errors.As` to detect wrapped errors.

Validation and input sanitization
- Use explicit DTO validation (e.g., `github.com/go-playground/validator/v10`) to reject malformed input early.

Concurrency & persistence
- The current database appears to be in-memory. For production or persistent needs:
  - Add thread-safety (mutex) to the in-memory store or use a concurrent data structure.
  - Replace/augment with a persistence layer (Postgres, SQLite, or a write-ahead log) and migrations.

Security and hardening
- Add request logging, request size limits, CORS config, and rate limiting where appropriate.
- Consider signing/hashing strategies and key management for production hashing schemes.

Operational
- Add a `Dockerfile` and Makefile targets for `docker build`/`docker run`.
- Add CI: GitHub Actions to run `go test ./...`, `go vet`, `golangci-lint`.
- Add code formatting/linting configs and pre-commit hooks.

API (current inferred endpoints)

- POST /transactions
  - Body: Transaction DTO JSON
  - Responses:
    - 201 Created — {"message":"Transaction created"}
    - 400 Bad Request — invalid JSON or validation error
    - <custom> — business errors returned by `TransactionError` (use its code and message)
    - 500 Internal Server Error — generic fallback

- GET /transactions
  - 200 OK — {"transactions": [...]}

- GET /transactions/validate
  - 200 OK — {"valid": true|false}

How to build & run (local)

Run with `go run` from repository root (zsh):

```zsh
# run app
go run ./cmd

# or build binary
go build -o immutable-ledger ./cmd
./immutable-ledger
```

Make targets
- If `Makefile` includes build/run/test targets, prefer `make` (check `Makefile`).

Suggested minimal roadmap (next tasks)

1. Fix HTTP handler response bug(s) in `app/transactions/transaction_routes.go` so handlers never write two responses for one request. Use `return` after writing a response.
2. Add unit tests for transaction service functions and handlers (happy + error paths).
3. Add concurrency safety or switch to persistent store (depending on goals).
4. Add CI that runs `go test ./...`, `go vet`, and a linter.

Notes & quick pointers
- Use `errors.As` when checking for typed errors if errors may be wrapped (fmt.Errorf with %w or third-party libs).
- Keep JSON response shapes consistent across endpoints to make client parsing predictable.

- ImmutableId prevent double transactions
- Final is about historical truth, the ledger will not revert any operation, it will only append a new transaction. A wrong transaction needs a a fix transaction to fix the mistake
- Validate the ledger will take as long its element size, so (O(n)), It`s a simple operation to validate. O(n) to get the account balance + time for fetching database data. Cache for accounts is acceptable, cache for save the state from X time of the ledger is acceptable too.