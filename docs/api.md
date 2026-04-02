# API Notes

Public document hosts:

- `GET /` returns exact stored plaintext bytes
- `HEAD /` returns headers only
- `ETag` is `"sha256-<content_sha256>"`
- `If-None-Match` returns `304 Not Modified` on exact match

Admin host routes under `/api/v1`:

- `GET /health`
- `GET /auth/check`
- `GET /documents`
- `GET /documents/:slug`
- `GET /documents/:slug/content`
- `POST /documents/:slug` create only
- `PUT /documents/:slug` replace only
- `DELETE /documents/:slug`

Authenticated admin requests require:

- `X-Timestamp`
- `X-Nonce`
- `X-Signature`

`GET /health` is public.

`GET /auth/check` is authenticated and returns `{"status":"ok"}` when signing succeeds.

Upload rules:

- `Content-Type` must be `text/plain` or `text/plain; charset=utf-8`
- body must be valid UTF-8
- max body size is controlled by `MAX_DOC_SIZE`
