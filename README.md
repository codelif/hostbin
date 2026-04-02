# hostbin

Host-routed plaintext document server with a Gin-based admin API and SQLite storage.

The repository root is the Go module root.

The current codebase is organized into:

- `cmd/server` for the server entrypoint
- `internal/domain` for reusable document, host, and slug logic
- `internal/protocol` for reusable API/auth protocol definitions
- `internal/server` for HTTP, auth middleware, storage, config, and runtime wiring

Quick starts:

- canonical install: `make build`, then `sudo make install-all`
- local development: `make test` or `make run`
- container workflow: `docker compose up --build`
