# Deployment

This page is the deployment index for `hostbin`.

Use it to choose the right deployment guide for your environment.

## Recommended production path

For a public VM deployment with wildcard hosts and Caddy in front, start here:

- `docs/deployment-cloudflare-caddy-systemd.md`

That guide covers DNS, TLS, Cloudflare, systemd, Caddy, verification, and basic operations.

## Deployment building blocks

- native service install: `docs/deployment-systemd.md`
- minimal Caddy reverse-proxy reference: `docs/deployment-caddy.md`
- minimal nginx reverse-proxy reference: `docs/deployment-nginx.md`

## Picking a guide

- evaluating locally -> `docs/getting-started.md`
- deploying on a Linux VM with `systemd` -> `docs/deployment-systemd.md`
- using Cloudflare and wildcard public hosts -> `docs/deployment-cloudflare-caddy-systemd.md`
- using Caddy but you already know how to manage certificates and services -> `docs/deployment-caddy.md`
- using nginx and managing TLS separately -> `docs/deployment-nginx.md`

## Invariants for every deployment

Regardless of proxy or init system:

- the app must receive the original `Host` header
- `ADMIN_HOST` must be an exact hostname under `BASE_DOMAIN`
- public documents are single-label subdomains only
- reverse proxy body limits must match `MAX_DOC_SIZE`
- if `TRUST_PROXY_HEADERS=true`, only trusted proxy CIDRs should be allowed to forward client IP/proto headers

## After deployment

Once the service is live:

- verify `GET /api/v1/health` on the admin host
- run `hbcli config check`
- create a test document
- fetch that document from its public hostname
- confirm logging, restart behavior, and backups

See `docs/operations.md` and `docs/troubleshooting.md` for day-2 guidance.
