# Architecture

The repository root contains the Go module, deployment notes, and both reusable and server-only packages.

Core runtime pieces:

- `cmd/server`: process entrypoint and graceful shutdown.
- `internal/domain/documents`: document types, repository interface, and domain errors.
- `internal/domain/hosts`: host normalization, host classification, and public URL building.
- `internal/domain/slugs`: shared slug validation rules.
- `internal/protocol/adminv1`: stable admin API DTOs, routes, and error codes.
- `internal/protocol/authsig`: canonical request signing primitives shared across transports.
- `internal/server/app`: composition root for storage, middleware, handlers, and server lifecycle.
- `internal/server/adminhttp`: admin Gin handlers and route registration.
- `internal/server/adminauth`: HMAC request verification middleware.
- `internal/server/publichttp`: public plaintext document handlers.
- `internal/server/dispatch`: strict host-based dispatch between public and admin HTTP stacks.
- `internal/server/store/sqlite`: SQLite schema initialization and document persistence.
- `internal/server/nonce`: atomic in-memory nonce replay protection.
- `internal/server/logging`: Zap logger creation and access logging middleware.

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
