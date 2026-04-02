# nginx

Reverse proxy requirements:

- preserve the original `Host` header
- terminate TLS before forwarding to the Go app
- do not rewrite request path, query string, or body
- keep the proxy body limit aligned with `MAX_DOC_SIZE`
- optionally IP-restrict the admin hostname at the proxy layer

Example server block:

```nginx
server {
    listen 443 ssl http2;
    server_name admin.domain.com *.domain.com;

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

Notes:

- use `$http_host` instead of `$host` so nginx does not normalize the incoming host value
- do not add a URI suffix to `proxy_pass`, or nginx may rewrite the upstream path
- if admin access should come only from a fixed location, add `allow` and `deny` rules on `admin.domain.com`

Smoke checks:

- public: `curl --resolve doc1.domain.com:443:127.0.0.1 https://doc1.domain.com/`
- admin health: `curl --resolve admin.domain.com:443:127.0.0.1 https://admin.domain.com/api/v1/health`
