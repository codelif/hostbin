# Cloudflare + Caddy + systemd

This is the recommended end-to-end production deployment guide for `hostbin` on a Linux VM.

It combines:

- Cloudflare DNS
- Caddy as the public TLS terminator and reverse proxy
- `systemd` for process management
- a native `hostbin` binary on localhost

Use this guide if you want one practical, real-world deployment path with wildcard public hosts such as `*.example.com`.

## What this covers

- DNS records
- Cloudflare SSL settings
- firewall expectations
- native server install
- `hostbin` environment configuration
- Caddy with the Cloudflare DNS challenge plugin
- service startup and verification

## Assumptions

- Linux host with `systemd`
- public IPv4 or IPv6 reachable on ports `80` and `443`
- `BASE_DOMAIN` such as `example.com`
- admin host such as `hbadmin.example.com`
- wildcard public hosts such as `hello.example.com`

## Why this path is recommended

`hostbin` needs the original `Host` header and often benefits from wildcard TLS for public document hosts. Caddy plus the Cloudflare DNS challenge covers that cleanly.

## 1. Prepare DNS in Cloudflare

Create these records for your zone:

- `A` or `AAAA` record for `hbadmin` -> your VM IP
- wildcard `A` or `AAAA` record for `*` -> your VM IP

You can keep them proxied in Cloudflare once the origin is working.

## 2. Set Cloudflare SSL mode

In Cloudflare SSL/TLS settings, use:

- `Full (strict)`

Do not use `Flexible`; the origin should present a valid certificate to Cloudflare.

## 3. Create a Cloudflare API token for DNS challenge

Create a token scoped to the zone with at least:

- `Zone:Read`
- `DNS:Edit`

Limit it to the specific zone for this deployment.

## 4. Open firewall ports

At the cloud firewall and host firewall layers, allow:

- `22/tcp` for SSH
- `80/tcp`
- `443/tcp`

Do not expose the app's local listen port publicly. Keep `hostbin` on `127.0.0.1:8080`.

If you use UFW:

```bash
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

## 5. Install build prerequisites

```bash
sudo apt update
sudo apt install -y build-essential git curl tar xz-utils sqlite3
```

Install a Go toolchain compatible with the repository. The project's Docker build can serve as the reference for the Go version in use.

## 6. Build and install `hostbin`

Create the service account:

```bash
sudo useradd --system --home /var/lib/hostbin --shell /usr/sbin/nologin hostbin
```

From the repository root:

```bash
make build
sudo make install-all
```

Fix ownership:

```bash
sudo mkdir -p /etc/hostbin /var/lib/hostbin
sudo chown -R hostbin:hostbin /etc/hostbin /var/lib/hostbin
sudo chmod 750 /var/lib/hostbin
```

## 7. Configure `hostbin`

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

Recommendations:

- reserve the admin label in `RESERVED_SUBDOMAINS`
- keep `LISTEN_ADDR` bound to localhost only
- set `TRUST_PROXY_HEADERS=true` only when traffic reaches the app through your local Caddy instance

Generate a secret if needed:

```bash
openssl rand -base64 48
```

Lock the file down:

```bash
sudo chown hostbin:hostbin /etc/hostbin/hostbin.env
sudo chmod 640 /etc/hostbin/hostbin.env
```

## 8. Start and verify `hostbin`

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now hostbin
sudo systemctl status hostbin --no-pager
```

Verify the app locally before adding Caddy:

```bash
curl -i -H 'Host: hbadmin.example.com' http://127.0.0.1:8080/api/v1/health
```

## 9. Build Caddy with the Cloudflare DNS plugin

Wildcard certificates require DNS challenge support. A stock Caddy package is often not enough for `*.example.com` unless it already includes the Cloudflare module.

Install `xcaddy`:

```bash
go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest
```

Build Caddy:

```bash
~/go/bin/xcaddy build --with github.com/caddy-dns/cloudflare
```

Install the binary:

```bash
sudo mv ./caddy /usr/local/bin/caddy
sudo chmod 755 /usr/local/bin/caddy
sudo setcap cap_net_bind_service=+ep /usr/local/bin/caddy
```

## 10. Create the Caddy service account and directories

```bash
sudo useradd --system --home /var/lib/caddy --shell /usr/sbin/nologin caddy
sudo mkdir -p /etc/caddy /var/lib/caddy
sudo chown -R caddy:caddy /etc/caddy /var/lib/caddy
```

## 11. Store the Cloudflare token for Caddy

Create `/etc/caddy/caddy.env`:

```dotenv
CLOUDFLARE_API_TOKEN=replace-with-your-cloudflare-api-token
```

Then secure it:

```bash
sudo chown root:caddy /etc/caddy/caddy.env
sudo chmod 640 /etc/caddy/caddy.env
```

## 12. Create the Caddyfile

Create `/etc/caddy/Caddyfile`:

```caddy
{
    email you@example.com
}

hbadmin.example.com, *.example.com {
    request_body {
        max_size 1MB
    }

    tls {
        dns cloudflare {env.CLOUDFLARE_API_TOKEN}
    }

    reverse_proxy 127.0.0.1:8080 {
        header_up Host {http.request.host}
        header_up X-Forwarded-For {client_ip}
        header_up X-Forwarded-Proto {scheme}
    }
}
```

Important details:

- `header_up Host {http.request.host}` is required for routing
- `max_size 1MB` should match `MAX_DOC_SIZE=1048576`
- the wildcard host and admin host can share one site block

## 13. Create a Caddy systemd unit

Create `/etc/systemd/system/caddy.service`:

```ini
[Unit]
Description=Caddy web server
Documentation=https://caddyserver.com/docs/
After=network-online.target
Wants=network-online.target

[Service]
Type=notify
User=caddy
Group=caddy
EnvironmentFile=/etc/caddy/caddy.env
ExecStart=/usr/local/bin/caddy run --environ --config /etc/caddy/Caddyfile
ExecReload=/usr/local/bin/caddy reload --config /etc/caddy/Caddyfile --force
TimeoutStopSec=5s
LimitNOFILE=1048576
PrivateTmp=true
AmbientCapabilities=CAP_NET_BIND_SERVICE
NoNewPrivileges=true

[Install]
WantedBy=multi-user.target
```

## 14. Validate and start Caddy

```bash
sudo /usr/local/bin/caddy validate --config /etc/caddy/Caddyfile
sudo systemctl daemon-reload
sudo systemctl enable --now caddy
sudo systemctl status caddy --no-pager
```

## 15. Verify the live deployment

Check admin health through the public hostname:

```bash
curl -i https://hbadmin.example.com/api/v1/health
```

Configure `hbcli`:

```bash
hbcli config init \
  --server-url https://hbadmin.example.com \
  --auth-key 'your-preshared-key'
```

Verify CLI connectivity:

```bash
hbcli config check
```

Create a test document:

```bash
printf 'hello from hostbin\n' | hbcli new hello
```

Fetch it publicly:

```bash
curl -i https://hello.example.com/
```

## Expected behavior

- `https://hbadmin.example.com/api/v1/health` -> `200` and `{"status":"ok"}`
- `hbcli config check` -> server reachable and authentication valid
- `https://hello.example.com/` -> exact plaintext document content
- `https://a.b.example.com/` -> not found
- `https://www.example.com/` -> not found if `www` is reserved

## Next steps

- native install details: `docs/deployment-systemd.md`
- operations: `docs/operations.md`
- troubleshooting: `docs/troubleshooting.md`
