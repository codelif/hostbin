# Caddy

Reverse proxy requirements:

- preserve the original `Host` header
- forward requests unchanged to the local app
- keep request body limits aligned with `MAX_DOC_SIZE`
- optionally restrict the admin host with matcher rules or a network firewall

Example Caddyfile:

```text
admin.domain.com, *.domain.com {
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

Notes:

- Caddy preserves the request URI by default; avoid adding handlers that rewrite the path or query string
- keep the `Host` header pass-through in place because host routing depends on it
- if you do not enable `TRUST_PROXY_HEADERS`, forwarded headers are ignored by the app except for raw logging fields

Smoke checks:

- public: `curl --resolve doc1.domain.com:443:127.0.0.1 https://doc1.domain.com/`
- admin health: `curl --resolve admin.domain.com:443:127.0.0.1 https://admin.domain.com/api/v1/health`
