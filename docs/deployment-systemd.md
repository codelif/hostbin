# systemd

This page covers native `hostbin` deployment under `systemd`.

Use this page if you want to run the server directly on a Linux host without containers.

## What this covers

- install layout
- service account
- environment file
- permissions
- service lifecycle
- verification and day-2 commands

## What this does not cover

- wildcard TLS
- DNS provider integration
- reverse proxy certificate management

For a full production example with Cloudflare and Caddy, use [Cloudflare + Caddy + systemd](deployment-cloudflare-caddy-systemd.md).

## Recommended install layout

- binary: `/usr/local/bin/hostbin`
- CLI: `/usr/local/bin/hbcli`
- env file: `/etc/hostbin/hostbin.env`
- state directory: `/var/lib/hostbin`
- database: `/var/lib/hostbin/data.db`
- unit file: `/etc/systemd/system/hostbin.service`

## 1. Create the service account

```bash
sudo useradd --system --home /var/lib/hostbin --shell /usr/sbin/nologin hostbin
```

## 2. Build and install

From the repository root:

```bash
make build
sudo make install-all
```

`make install-all`:

- installs `hostbin` to `/usr/local/bin/hostbin`
- installs `hbcli` to `/usr/local/bin/hbcli`
- installs the unit to `/etc/systemd/system/hostbin.service`
- installs `/etc/hostbin/hostbin.env` only if it does not already exist
- creates `/var/lib/hostbin`

## 3. Configure the environment file

Edit `/etc/hostbin/hostbin.env`.

Example:

```dotenv
LISTEN_ADDR=127.0.0.1:8080
BASE_DOMAIN=example.com
ADMIN_HOST=hbadmin.example.com
PRESHARED_KEY=replace-with-a-long-random-secret-at-least-32-bytes
DB_PATH=/var/lib/hostbin/data.db
RESERVED_SUBDOMAINS=hbadmin,www,api
MAX_DOC_SIZE=1048576
AUTH_TIMESTAMP_SKEW_SECONDS=60
NONCE_TTL_SECONDS=300
TRUST_PROXY_HEADERS=true
TRUSTED_PROXY_CIDRS=127.0.0.1/32,::1/128
LOG_LEVEL=info
```

Field notes:

- `LISTEN_ADDR` should usually be `127.0.0.1:8080` behind a reverse proxy
- `ADMIN_HOST` must be under `BASE_DOMAIN`
- `PRESHARED_KEY` must be at least 32 bytes
- `RESERVED_SUBDOMAINS` should include the admin label and any labels you never want claimable as public documents
- enable `TRUST_PROXY_HEADERS` only when requests arrive from a trusted proxy; see [Overview: Reverse proxy requirements](overview.md#reverse-proxy-requirements)

Generate a secret with:

```bash
openssl rand -base64 48
```

## 4. Set ownership and permissions

```bash
sudo mkdir -p /etc/hostbin /var/lib/hostbin
sudo chown -R hostbin:hostbin /etc/hostbin /var/lib/hostbin
sudo chmod 640 /etc/hostbin/hostbin.env
sudo chmod 750 /var/lib/hostbin
```

Keep `PRESHARED_KEY` only in the env file, not in the unit file.

## 5. Enable and start the service

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now hostbin
```

## 6. Verify the unit

Validate the unit file:

```bash
sudo systemd-analyze verify /etc/systemd/system/hostbin.service
```

Check service state:

```bash
sudo systemctl status hostbin --no-pager
```

Verify the app locally before introducing the reverse proxy:

```bash
curl -i -H 'Host: hbadmin.example.com' http://127.0.0.1:8080/api/v1/health
```

Expected response:

```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{"status":"ok"}
```

## Reverse proxy reminder

Pair this with a reverse proxy that:

- preserves the original `Host` header
- does not rewrite path or query string
- keeps body limits aligned with `MAX_DOC_SIZE`

See:

- [Overview: Reverse proxy requirements](overview.md#reverse-proxy-requirements)
- [Caddy](deployment-caddy.md)
- [nginx](deployment-nginx.md)
- [Cloudflare + Caddy + systemd](deployment-cloudflare-caddy-systemd.md)

## Day-2 commands

Follow logs:

```bash
sudo journalctl -u hostbin -f
```

Restart after config changes:

```bash
sudo systemctl restart hostbin
```

Stop and start manually:

```bash
sudo systemctl stop hostbin
sudo systemctl start hostbin
```

## Next steps

- full Cloudflare + Caddy production runbook: [Cloudflare + Caddy + systemd](deployment-cloudflare-caddy-systemd.md)
- operations: [Operations](operations.md)
- troubleshooting: [Troubleshooting](troubleshooting.md)
