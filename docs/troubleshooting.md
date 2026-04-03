# Troubleshooting

This page lists common `hostbin` deployment and usage failures with likely causes and quick verification steps.

Use it when the service is reachable but behaving unexpectedly, or when local development and deployment checks do not match the documented happy path.

## Admin health endpoint returns 404

Likely causes:

- the request did not reach the exact `ADMIN_HOST`
- the reverse proxy rewrote `Host`
- `ADMIN_HOST` is misconfigured

Verify:

```bash
curl -i -H 'Host: hbadmin.example.com' http://127.0.0.1:8080/api/v1/health
```

Fix:

- confirm `ADMIN_HOST` matches the hostname you are using
- confirm the proxy passes through the original `Host` header; see [Overview: Reverse proxy requirements](overview.md#reverse-proxy-requirements), [Caddy](deployment-caddy.md), and [nginx](deployment-nginx.md)

## Public hostname returns 404

Likely causes:

- document does not exist
- slug is reserved
- hostname has more than one subdomain label
- request is not under `BASE_DOMAIN`

Verify:

- `hbcli list`
- `hbcli info <slug>`

Fix:

- create the missing document
- use a non-reserved slug
- use `slug.example.com`, not `a.b.example.com`; the routing rules are summarized in [Overview: Routing model](overview.md#routing-model)

## `hbcli config check` fails health check

Likely causes:

- `server_url` points to the wrong host
- DNS or `/etc/hosts` does not resolve the admin hostname
- the service is not listening where expected

Verify:

```bash
curl -i https://hbadmin.example.com/api/v1/health
sudo systemctl status hostbin --no-pager
```

Fix:

- set `server_url` to the exact admin host
- fix DNS or your local hosts entry
- restart the service if needed; see [CLI: Bootstrap config](cli.md#bootstrap-config) and [Operations](operations.md#service-control)

## Local development works with curl but not comfortably with `hbcli`

Likely cause:

- you are trying to use plain `localhost` even though the app routes by subdomain hostnames

Fix:

- use `BASE_DOMAIN=lvh.me`
- use `ADMIN_HOST=hbadmin.lvh.me`
- point `hbcli` at `http://hbadmin.lvh.me:8080`; the full example lives in [Getting Started](getting-started.md#6-initialize-hbcli)

Example:

```bash
hbcli config init \
  --server-url http://hbadmin.lvh.me:8080 \
  --auth-key 'your-preshared-key'
```

Fallback:

- for debugging only, use manual `Host` headers with `curl`

## `hbcli config check` fails auth check

Likely causes:

- `auth_key` does not match `PRESHARED_KEY`
- request path or host is being rewritten by the proxy
- system clock drift is too large

Fix:

- update the CLI config with the correct key
- confirm the proxy preserves `Host`, path, and query string
- sync server and client clocks; see [CLI](cli.md) and [API: Admin auth model](api.md#admin-auth-model)

## Uploads fail with document too large

Likely causes:

- request exceeded `MAX_DOC_SIZE`
- reverse proxy body limit is lower than the app limit

Fix:

- increase `MAX_DOC_SIZE` if appropriate
- raise the proxy body limit to the same value; see [Caddy](deployment-caddy.md) and [nginx](deployment-nginx.md)

## Uploads fail with `bad_request`

Likely causes:

- wrong `Content-Type`
- malformed request body

Fix:

- send `text/plain` or `text/plain; charset=utf-8`
- use `hbcli` if you are unsure about request formatting; the request rules are listed in [API: Upload requirements](api.md#upload-requirements)

## Uploads fail with `invalid_utf8`

Cause:

- the body is not valid UTF-8

Fix:

- re-encode the source file as UTF-8 before upload; see [API: Upload requirements](api.md#upload-requirements)

## Wildcard TLS certificate issuance fails

Likely causes:

- Caddy was built without the Cloudflare DNS module
- Cloudflare API token is missing or lacks `Zone:Read` or `DNS:Edit`
- the token is not available to the Caddy service environment

Verify:

```bash
sudo journalctl -u caddy -n 100 --no-pager
```

Fix:

- rebuild Caddy with `github.com/caddy-dns/cloudflare`
- verify token scope and environment file permissions; see the [Caddy DNS challenge docs](https://caddyserver.com/docs/automatic-https#dns-challenge) and [Cloudflare API token docs](https://developers.cloudflare.com/fundamentals/api/get-started/create-token/)

## Cloudflare is enabled but HTTPS still fails

Likely causes:

- Cloudflare SSL mode is not `Full (strict)`
- origin certificate issuance has not completed

Fix:

- set Cloudflare SSL/TLS mode to `Full (strict)`; see the [Cloudflare encryption mode docs](https://developers.cloudflare.com/ssl/origin-configuration/ssl-modes/full-strict/)
- verify Caddy is running and has obtained certificates successfully

## Client IPs or scheme look wrong in logs

Likely causes:

- `TRUST_PROXY_HEADERS=false`
- trusted proxy CIDRs do not match the proxy address

Fix:

- enable `TRUST_PROXY_HEADERS` only when traffic comes through a trusted proxy
- set `TRUSTED_PROXY_CIDRS` correctly for that proxy path; see [Overview: Reverse proxy requirements](overview.md#reverse-proxy-requirements)

## Delete fails in automation

Cause:

- `hbcli delete` requires `--yes` in non-interactive mode

Fix:

```bash
hbcli delete hello --yes
```

## Need deeper visibility

Check logs:

```bash
sudo journalctl -u hostbin -f
sudo journalctl -u caddy -f
```

Then run the smoke checks from [Operations](operations.md#suggested-smoke-checks-after-any-change).

## See also

- [Operations](operations.md) for service control, logs, and smoke checks
- [Deployment](deployment.md) for picking the right deployment reference
- [CLI](cli.md) for config and command usage details
