# Getting Started

This page walks through a local, production-like run using the native server and `hbcli`.

Use this page if you want to prove the project works on your machine before setting up production DNS, TLS, or a reverse proxy.

## What this covers

- building `hostbin` and `hbcli`
- starting the server locally
- configuring `hbcli`
- creating and fetching a document

## What this does not cover

- production TLS
- wildcard certificates
- reverse proxy setup
- hardened systemd deployment

For those topics, jump to [Deployment](deployment.md) or the full [Cloudflare + Caddy + systemd guide](deployment-cloudflare-caddy-systemd.md).

## Prerequisites

- Go toolchain compatible with the repository
- `make`
- outbound DNS resolution for `lvh.me`

## 1. Build the binaries

From the repository root:

```bash
make build
```

This creates:

- `bin/hostbin`
- `bin/hbcli`

## 2. Create a local config file

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env` so it uses a wildcard loopback domain. This example uses `lvh.me`, which resolves wildcard subdomains to `127.0.0.1`:

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

Generate a secret if you need one:

```bash
openssl rand -base64 48
```

## 3. Why `lvh.me` is the easiest local setup

`hostbin` routes requests by hostname, not by path. Plain `localhost` does not give you realistic subdomains for either the admin host or public documents.

Using `lvh.me` avoids that problem:

- `hbadmin.lvh.me` resolves to `127.0.0.1`
- `hello.lvh.me` resolves to `127.0.0.1`
- `hbcli` can talk to the exact admin host without custom headers
- public document fetches work without manual `Host` overrides

If you need a fully offline setup, see the [fallback note](#10-fallback-if-lvhme-is-not-suitable) at the end of this page and the related [troubleshooting entry](troubleshooting.md#local-development-works-with-curl-but-not-comfortably-with-hbcli).

## 4. Start the server

```bash
make run
```

The server loads values from `.env` if present and listens on `127.0.0.1:8080`.

## 5. Verify the admin health endpoint

In another terminal:

```bash
curl -i http://hbadmin.lvh.me:8080/api/v1/health
```

Expected response:

```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{"status":"ok"}
```

## 6. Initialize `hbcli`

Create local CLI config:

```bash
./bin/hbcli config init \
  --server-url http://hbadmin.lvh.me:8080 \
  --auth-key 'your-preshared-key'
```

Then verify connectivity and auth:

```bash
./bin/hbcli config check
```

Expected output:

```text
OK server reachable
OK authentication valid
```

The config file is written to `~/.config/hbcli/config.toml` unless you pass `--config`. For more bootstrap options, see [CLI: Bootstrap config](cli.md#bootstrap-config).

## 7. Create a document

```bash
printf 'hello from hostbin\n' | ./bin/hbcli new hello
```

This creates slug `hello`, which maps to public hostname `hello.lvh.me`.

## 8. Fetch the public document locally

With `lvh.me`, you can fetch the public document directly without overriding headers:

```bash
curl -i http://hello.lvh.me:8080/
```

Expected response body:

```text
hello from hostbin
```

## 9. Try a few more CLI commands

List documents:

```bash
./bin/hbcli list
```

Show metadata:

```bash
./bin/hbcli info hello
```

Print content:

```bash
./bin/hbcli get hello
```

Edit in your configured editor:

```bash
./bin/hbcli edit hello
```

Delete the document:

```bash
./bin/hbcli delete hello
```

## 10. Fallback if `lvh.me` is not suitable

If your environment cannot rely on public wildcard DNS helpers:

- use local DNS such as `dnsmasq` or `CoreDNS`, or
- add only the admin host to `/etc/hosts` and use manual `Host` headers for public fetches

Example fallback public test:

```bash
curl -i -H 'Host: hello.example.test' http://127.0.0.1:8080/
```

This fallback is useful for debugging, but it is not the recommended default onboarding path. If hostname resolution is the main issue, see [Troubleshooting](troubleshooting.md#local-development-works-with-curl-but-not-comfortably-with-hbcli).

## 11. Clean up

- stop `make run`
- remove `data.db` if you want a clean state

## Next steps

- CLI details: [CLI](cli.md)
- daily CLI commands: [CLI: Daily commands](cli.md#daily-commands)
- full deployment options: [Deployment](deployment.md)
- recommended production guide: [Cloudflare + Caddy + systemd](deployment-cloudflare-caddy-systemd.md)
