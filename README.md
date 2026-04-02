# hostbin

Host-routed plaintext document server with a Gin-based admin API and SQLite storage.

The repository root is the Go module root.

The current codebase is organized into:

- `cmd/server` for the server entrypoint
- `cmd/hbcli` for the CLI entrypoint
- `internal/domain` for reusable document, host, and slug logic
- `internal/protocol` for reusable API/auth protocol definitions
- `internal/cli` for CLI config, prompting, HTTP client, and commands
- `internal/server` for HTTP, auth middleware, storage, config, and runtime wiring

Quick starts:

- canonical install: `make build`, then `sudo make install-all`
- first-time CLI setup: `hbcli config init`
- local development: copy `.env.example` to `.env`, then use `make test` or `make run`
- container workflow: `docker compose up --build`
