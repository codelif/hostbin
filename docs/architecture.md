# Architecture

The repository root contains documentation and deployment notes. The runnable server module lives in `server/`.

Core runtime pieces:

- `cmd/server`: process entrypoint and graceful shutdown.
- `internal/app`: wiring for config, logging, storage, auth, Gin engines, and top-level HTTP wrappers.
- `internal/router`: strict host normalization and dispatch between admin and public handlers.
- `internal/public`: read-only plaintext document serving with ETag support.
- `internal/admin`: authenticated JSON admin API for CRUD operations.
- `internal/auth`: HMAC request verification and canonical request handling.
- `internal/storage/sqlite`: SQLite schema initialization and document persistence.
- `internal/nonce`: atomic in-memory nonce replay protection.
- `internal/logging`: Zap logger and access logging.

Phase-by-phase implementation order:

1. bootstrap repository, config, slug validation, and host routing
2. add SQLite schema and document store
3. add public plaintext serving
4. add admin CRUD handlers
5. add HMAC verification and nonce replay protection
6. add request IDs, Zap logs, and proxy trust rules
7. add tests and deployment docs

Suggested commit plan:

1. `chore: bootstrap hostbin repository`
2. `feat: add config loading and host routing`
3. `feat: add sqlite-backed document storage`
4. `feat: serve public documents by hostname`
5. `feat: add admin document api`
6. `feat: add hmac auth and atomic nonce protection`
7. `feat: enforce utf-8 uploads and add structured logging`
8. `test: add unit and integration coverage`
9. `docs: add deployment and architecture notes`
