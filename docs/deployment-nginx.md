# nginx

This page is a minimal nginx reverse-proxy reference for `hostbin`.

It is intentionally not a full deployment guide. It focuses on the nginx settings the application requires.

## Reverse proxy requirements

- preserve the original `Host` header; see [API: Host routing requirements](api.md#host-routing-requirements)
- terminate TLS before forwarding to the Go app
- do not rewrite request path, query string, or body
- keep proxy body limits aligned with `MAX_DOC_SIZE`
- optionally IP-restrict the admin hostname at the proxy layer

## Minimal server block

```nginx
server {
    listen 443 ssl http2;
    server_name hbadmin.example.com *.example.com;

    client_max_body_size 1m;

    location / {
        proxy_http_version 1.1;
        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_pass http://127.0.0.1:8080;
    }
}
```

## Why this matters

- use `$http_host` instead of `$host` so nginx does not normalize the incoming host value
- do not add a URI suffix to `proxy_pass`, or nginx may rewrite the upstream path
- `client_max_body_size` must match the app-side upload size limit

The upstream directives are documented in the nginx docs for [`proxy_set_header`](https://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_set_header), [`proxy_pass`](https://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_pass), [`client_max_body_size`](https://nginx.org/en/docs/http/ngx_http_core_module.html#client_max_body_size), and [`server_name`](https://nginx.org/en/docs/http/server_names.html).

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

- native service deployment: [systemd](deployment-systemd.md)
- production operations: [Operations](operations.md)
- troubleshooting: [Troubleshooting](troubleshooting.md)
