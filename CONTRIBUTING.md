# Contributing

Thanks for your interest in contributing to `hostbin`.

This project is small and opinionated. The best contributions keep the routing model clear, the deployment story practical, and the CLI and server behavior predictable.

## Development setup

Prerequisites:

- Go toolchain version compatible with `go.mod`
- `make`

Common commands:

```bash
make build
make test
make run
```

Built binaries are written to `bin/`.

## Local development

For local development, you can either use a wildcard loopback domain such as `lvh.me` or edit your `/etc/hosts` file. The [Getting Started guide](docs/getting-started.md) explains the recommended local flow, and [Why `lvh.me` is the easiest local setup](docs/getting-started.md#3-why-lvhme-is-the-easiest-local-setup) covers why hostname-based routing matters.

Example `.env` values:

```dotenv
LISTEN_ADDR=127.0.0.1:8080
BASE_DOMAIN=lvh.me
ADMIN_HOST=hbadmin.lvh.me
PRESHARED_KEY=replace-with-a-long-random-secret-at-least-32-bytes
DB_PATH=./data.db
RESERVED_SUBDOMAINS=hbadmin,www,api
MAX_DOC_SIZE=1048576
AUTH_TIMESTAMP_SKEW_SECONDS=60
NONCE_TTL_SECONDS=300
TRUST_PROXY_HEADERS=false
TRUSTED_PROXY_CIDRS=127.0.0.1/32,::1/128
LOG_LEVEL=info
```

Why `lvh.me`:

- the server routes on `Host`
- `hbcli` needs to talk to the real admin hostname
- public documents also live on subdomains
- saves you from manually crafting curl requests

Typical local flow:

```bash
cp .env.example .env
make run
./bin/hbcli config init --server-url http://hbadmin.lvh.me:8080 --auth-key 'your-preshared-key'
./bin/hbcli config check
printf 'hello from hostbin\n' | ./bin/hbcli new hello
curl http://hello.lvh.me:8080/
```

## Project structure

- `cmd/server` - server entrypoint
- `cmd/hbcli` - CLI entrypoint
- `internal/domain` - domain rules for hosts, slugs, and documents
- `internal/protocol` - stable API and auth protocol definitions
- `internal/cli` - CLI configuration, commands, editor support, and HTTP client
- `internal/server` - app wiring, auth, dispatch, handlers, logging, and storage
- `docs` - product, deployment, operations, and troubleshooting docs

See also:

- [README](README.md) for the top-level project map and development commands
- [Architecture](docs/architecture.md) for contributor-oriented package layout notes
- [Getting Started](docs/getting-started.md) for the local evaluation workflow

## Expectations for changes

- keep the host-based routing model intact
- preserve the exact admin-host versus public-host split
- avoid proxy-dependent assumptions in core request handling
- update documentation when behavior changes
- add or update tests for non-trivial logic changes

## Before opening a pull request

Please run:

```bash
make test
```

If your change affects documentation, examples, or local workflows, sanity-check the relevant docs as well. The [development commands in `README.md`](README.md#development-commands) are the baseline checks.

## Pull request guidance

- explain the user-visible reason for the change
- mention any deployment or compatibility impact
- keep diffs focused when possible
- include docs updates for behavior, config, or workflow changes
