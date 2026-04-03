# Overview

This page explains what `hostbin` is, how requests are routed, and what assumptions matter in real deployments.

Use this page if you are deciding whether the project fits your use case or you want to understand the routing model before reading deployment docs.

## What hostbin is

`hostbin` stores plaintext documents keyed by slug and serves each document from its own hostname.

For a base domain such as `example.com`:

- the admin API lives on one exact hostname such as `hbadmin.example.com`
- public documents live on `https://<slug>.example.com/`
- the application decides whether a request is public or admin by classifying the incoming `Host` header

This means the reverse proxy is part of the runtime contract. If it rewrites or normalizes `Host` incorrectly, routing breaks.

For local development, prefer a wildcard loopback domain such as `lvh.me` instead of `localhost` so the app can still receive realistic subdomain hosts.

## Core concepts

- `BASE_DOMAIN` - parent domain for public documents, such as `example.com`
- `ADMIN_HOST` - exact admin hostname, such as `hbadmin.example.com`
- `slug` - the document identifier and public subdomain label
- `RESERVED_SUBDOMAINS` - labels that may never be used as public slugs

## Routing model

Given:

- `BASE_DOMAIN=example.com`
- `ADMIN_HOST=hbadmin.example.com`
- `RESERVED_SUBDOMAINS=hbadmin,www,api`

The runtime behavior is:

- `hbadmin.example.com` -> admin API
- `hello.example.com` -> public document for slug `hello`
- `www.example.com` -> rejected because `www` is reserved
- `api.example.com` -> rejected because `api` is reserved
- `a.b.example.com` -> rejected because public hosts must be a single subdomain label under `BASE_DOMAIN`
- `example.com` -> rejected because the app expects either the exact admin host or a public slug host

## Public document behavior

Public hosts respond on `/`:

- `GET /` returns the exact stored plaintext bytes
- `HEAD /` returns headers only
- `ETag` is returned as `"sha256-<content_sha256>"`
- `If-None-Match` on an exact match returns `304 Not Modified`

The public surface is intentionally narrow. It is designed to serve the stored plaintext document for a host, not a browsable app.

## Admin API behavior

Admin routes live under `/api/v1` on the exact `ADMIN_HOST`.

Key endpoints:

- `GET /api/v1/health` - public health endpoint
- `GET /api/v1/auth/check` - signed auth verification
- `GET /api/v1/documents` - list documents
- `POST /api/v1/documents/:slug` - create document
- `PUT /api/v1/documents/:slug` - replace document
- `DELETE /api/v1/documents/:slug` - delete document

Authenticated admin requests require signed headers:

- `X-Timestamp`
- `X-Nonce`
- `X-Signature`

The bundled `hbcli` client handles signing for you.

## Upload constraints

- `Content-Type` must be `text/plain` or `text/plain; charset=utf-8`
- request body must be valid UTF-8
- maximum upload size is controlled by `MAX_DOC_SIZE`
- the reverse proxy body limit must match `MAX_DOC_SIZE`

## Storage model

- metadata and document content are stored in SQLite
- default state location in systemd deployments is `/var/lib/hostbin/data.db`
- the server creates schema as needed when the database is opened

## Reverse proxy requirements

Any reverse proxy in front of `hostbin` must:

- preserve the original `Host` header
- forward the request path and query string unchanged
- keep body size limits aligned with `MAX_DOC_SIZE`
- terminate TLS before forwarding to the local app

If you enable `TRUST_PROXY_HEADERS=true`, only trusted proxy IP ranges should be allowed to set forwarded headers.

## Recommended reading order

- local evaluation: `docs/getting-started.md`
- CLI workflows: `docs/cli.md`
- full production deployment: `docs/deployment-cloudflare-caddy-systemd.md`
- operational guidance: `docs/operations.md`
