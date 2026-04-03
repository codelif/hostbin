# Caddy

This page is a minimal Caddy reverse-proxy reference for `hostbin`.

It is intentionally not a full deployment guide. It focuses on the proxy behavior the application requires.

If you need wildcard TLS, Cloudflare integration, service setup, and end-to-end host deployment, use `docs/deployment-cloudflare-caddy-systemd.md` instead.

## Reverse proxy requirements

- preserve the original `Host` header
- forward requests to the local app without rewriting path or query string
- keep request body limits aligned with `MAX_DOC_SIZE`
- optionally restrict the admin host separately with Caddy matchers or a firewall

## Minimal Caddyfile example

```caddy
hbadmin.example.com, *.example.com {
    request_body {
        max_size 1MB
    }

    reverse_proxy 127.0.0.1:8080 {
        header_up Host {http.request.host}
        header_up X-Forwarded-For {client_ip}
        header_up X-Forwarded-Proto {scheme}
    }
}
```

## Why this matters

- `Host` pass-through is required because host routing happens inside the application
- path and query rewriting can break auth signatures and route matching
- body limits that are lower than `MAX_DOC_SIZE` cause uploads to fail at the proxy layer

## Notes

- Caddy preserves the request URI by default; avoid adding handlers that rewrite the path or query string
- if you do not enable `TRUST_PROXY_HEADERS`, forwarded headers are ignored by the app except for raw logging fields
- if you do enable `TRUST_PROXY_HEADERS`, restrict `TRUSTED_PROXY_CIDRS` to the proxy IP ranges that actually forward traffic

## Smoke checks

Public host:

```bash
curl --resolve doc1.example.com:443:127.0.0.1 https://doc1.example.com/
```

Admin health:

```bash
curl --resolve hbadmin.example.com:443:127.0.0.1 https://hbadmin.example.com/api/v1/health
```

## Related docs

- native service deployment: `docs/deployment-systemd.md`
- full Cloudflare + Caddy runbook: `docs/deployment-cloudflare-caddy-systemd.md`
- troubleshooting: `docs/troubleshooting.md`
