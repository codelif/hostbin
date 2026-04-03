# hostbin

Host-routed plaintext document server

`hostbin` is opinionated on purpose:

- public documents live at `https://<slug>.<base-domain>/`
- admin API traffic lives on one exact hostname such as `https://hbadmin.<base-domain>/api/v1/...`
- uploads are plaintext only and must be valid UTF-8
- admin writes are authenticated with signed request headers
- one binary serves both the public and admin surfaces by dispatching on `Host`

## Highlights

- single binary server with SQLite persistence
- strict host-based routing between public and admin traffic
- plaintext reads with `ETag` and `If-None-Match` support
- HMAC-signed admin API with nonce replay protection
- `hbcli` for config, upload, edit, list, and delete workflows
- straightforward reverse-proxy deployment behind Caddy or nginx

## How Routing Works

Given:

- `BASE_DOMAIN=example.com`
- `ADMIN_HOST=hbadmin.example.com`

Then:

- `https://hbadmin.example.com/api/v1/health` -> admin API
- `https://hello.example.com/` -> plaintext content for slug `hello`
- `https://www.example.com/` -> rejected if `www` is reserved
- `https://a.b.example.com/` -> rejected; public hosts must be single-label subdomains

This routing model is enforced by the application, not just by the reverse proxy. The proxy must preserve the original `Host` header.

## Quickstart

For a local evaluation flow, see `docs/getting-started.md`.

High level:

1. copy `.env.example` to `.env`
2. use a wildcard loopback domain such as `lvh.me` for local testing
3. run `make run`
4. configure `hbcli`
5. create a document and fetch it through its public hostname

For local development, prefer:

- `BASE_DOMAIN=lvh.me`
- `ADMIN_HOST=hbadmin.lvh.me`

That lets `hbcli` talk to `http://hbadmin.lvh.me:8080` and lets you fetch public docs at `http://hello.lvh.me:8080/` without custom `Host` headers.

Otherwise you can edit your /etc/hosts file

## Deployment Paths

- recommended production runbook: `docs/deployment-cloudflare-caddy-systemd.md`
- deployment index: `docs/deployment.md`
- native service install details: `docs/deployment-systemd.md`
- minimal Caddy proxy reference: `docs/deployment-caddy.md`
- minimal nginx proxy reference: `docs/deployment-nginx.md`

## Documentation

- product and routing model: `docs/overview.md`
- local first run: `docs/getting-started.md`
- CLI usage: `docs/cli.md`
- API reference: `docs/api.md`
- production operations: `docs/operations.md`
- troubleshooting: `docs/troubleshooting.md`
- architecture for contributors: `docs/architecture.md`

## Repository Layout

- `cmd/server` - server entrypoint
- `cmd/hbcli` - CLI entrypoint
- `internal/domain` - host, slug, and document domain logic
- `internal/protocol` - API and auth protocol definitions
- `internal/cli` - CLI config, prompting, editor, and HTTP client
- `internal/server` - config, auth, dispatch, HTTP handlers, storage, and runtime wiring
- `deploy/systemd` - packaged service and env examples
- `docs` - user, operator, and contributor documentation

## Development Commands

- `make build` - build `hostbin` and `hbcli` into `./bin`
- `make test` - run the Go test suite
- `make run` - run the server with environment from `.env` if present
- `docker compose up --build` - containerized local workflow

## First Things To Read

- evaluating the project: `docs/overview.md`
- trying it locally: `docs/getting-started.md`
- deploying it on a VM: `docs/deployment-cloudflare-caddy-systemd.md`
- integrating with the admin API: `docs/api.md`
