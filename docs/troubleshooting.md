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
- confirm the proxy passes through the original `Host` header

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
- use `slug.example.com`, not `a.b.example.com`

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
- restart the service if needed

## Local development works with curl but not comfortably with `hbcli`

Likely cause:

- you are trying to use plain `localhost` even though the app routes by subdomain hostnames

Fix:

- use `BASE_DOMAIN=lvh.me`
- use `ADMIN_HOST=hbadmin.lvh.me`
- point `hbcli` at `http://hbadmin.lvh.me:8080`

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
- sync server and client clocks

## Uploads fail with document too large

Likely causes:

- request exceeded `MAX_DOC_SIZE`
- reverse proxy body limit is lower than the app limit

Fix:

- increase `MAX_DOC_SIZE` if appropriate
- raise the proxy body limit to the same value

## Uploads fail with `bad_request`

Likely causes:

- wrong `Content-Type`
- malformed request body

Fix:

- send `text/plain` or `text/plain; charset=utf-8`
- use `hbcli` if you are unsure about request formatting

## Uploads fail with `invalid_utf8`

Cause:

- the body is not valid UTF-8

Fix:

- re-encode the source file as UTF-8 before upload

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
- verify token scope and environment file permissions

## Cloudflare is enabled but HTTPS still fails

Likely causes:

- Cloudflare SSL mode is not `Full (strict)`
- origin certificate issuance has not completed

Fix:

- set Cloudflare SSL/TLS mode to `Full (strict)`
- verify Caddy is running and has obtained certificates successfully

## Client IPs or scheme look wrong in logs

Likely causes:

- `TRUST_PROXY_HEADERS=false`
- trusted proxy CIDRs do not match the proxy address

Fix:

- enable `TRUST_PROXY_HEADERS` only when traffic comes through a trusted proxy
- set `TRUSTED_PROXY_CIDRS` correctly for that proxy path

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

Then run smoke checks from `docs/operations.md`.
