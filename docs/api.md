# API

This page documents the HTTP surface exposed by `hostbin`.

The server has two distinct HTTP faces:

- public document reads on public slug hosts
- admin API routes on the exact admin host under `/api/v1`

## Host routing requirements

The application classifies requests by `Host`.

Given:

- `BASE_DOMAIN=example.com`
- `ADMIN_HOST=hbadmin.example.com`

Then:

- `hbadmin.example.com` -> admin API
- `hello.example.com` -> public document for slug `hello`
- `a.b.example.com` -> invalid

If a reverse proxy rewrites `Host`, the request may be routed to the wrong surface or rejected entirely.

## Public document endpoints

Public hosts serve the stored plaintext document on `/`.

### `GET /`

Returns the exact stored plaintext bytes.

Headers:

- `Content-Type: text/plain; charset=utf-8`
- `ETag: "sha256-<content_sha256>"`

### `HEAD /`

Returns the same headers as `GET /` without the response body.

### Conditional requests

If `If-None-Match` matches the current `ETag`, the server returns `304 Not Modified`.

Example:

```bash
curl -i https://hello.example.com/
curl -i -H 'If-None-Match: "sha256-<sha256>"' https://hello.example.com/
```

## Admin API base path

All admin routes live under:

```text
/api/v1
```

## Admin auth model

All admin routes except `GET /api/v1/health` require signed headers:

- `X-Timestamp`
- `X-Nonce`
- `X-Signature`

The bundled `hbcli` client signs requests for you. If you are writing your own client, the signature input is built from:

```text
<METHOD>
<RAW_PATH>
<RAW_QUERY>
<SHA256_HEX_OF_BODY>
<UNIX_TIMESTAMP>
<NONCE>
```

`X-Signature` is the hex-encoded HMAC-SHA256 of that canonical string using the shared secret.

The shared secret must be at least 32 bytes.

## Admin routes

### `GET /api/v1/health`

Public health endpoint.

Response:

```json
{"status":"ok"}
```

### `GET /api/v1/auth/check`

Authenticated auth validation endpoint.

Response on success:

```json
{"status":"ok"}
```

### `GET /api/v1/documents`

List document metadata.

Response:

```json
{
  "documents": [
    {
      "slug": "hello",
      "url": "https://hello.example.com/",
      "size_bytes": 19,
      "sha256": "...",
      "created_at": "2026-04-03T12:00:00Z",
      "updated_at": "2026-04-03T12:00:00Z"
    }
  ]
}
```

### `GET /api/v1/documents/:slug`

Fetch one document's metadata.

### `GET /api/v1/documents/:slug/content`

Fetch one document's raw plaintext content.

Returns:

- `200 OK`
- `Content-Type: text/plain; charset=utf-8`

### `POST /api/v1/documents/:slug`

Create a new document. Fails with `409` if the document already exists.

### `PUT /api/v1/documents/:slug`

Replace an existing document. Fails with `404` if the document does not exist.

### `DELETE /api/v1/documents/:slug`

Delete an existing document.

Success response:

```json
{
  "deleted": true,
  "slug": "hello"
}
```

## Upload requirements

For `POST` and `PUT` document writes:

- `Content-Type` must be `text/plain` or `text/plain; charset=utf-8`
- body must be valid UTF-8
- body size must not exceed `MAX_DOC_SIZE`

Example content upload with a custom client after signing:

```http
POST /api/v1/documents/hello HTTP/1.1
Host: hbadmin.example.com
Content-Type: text/plain; charset=utf-8
X-Timestamp: 1712145600
X-Nonce: 0123456789abcdef0123456789abcdef
X-Signature: <hex-hmac>

hello from hostbin
```

## Error codes

JSON error responses use this shape:

```json
{"error":"not_found"}
```

Known error codes:

- `unauthorized`
- `invalid_signature`
- `invalid_timestamp`
- `replayed_nonce`
- `already_exists`
- `invalid_slug`
- `not_found`
- `document_too_large`
- `bad_request`
- `invalid_utf8`
- `method_not_allowed`
- `internal_error`

Typical HTTP status mappings:

- `400` - invalid slug, bad content type, invalid UTF-8
- `401` - auth/signature/timestamp/nonce failures
- `404` - missing document or invalid host routing target
- `409` - create on an existing document
- `413` - document too large
- `500` - unexpected internal failure

## Recommended client path

For shell use and most automation, prefer `hbcli` over hand-rolled signing logic.

See `docs/cli.md` for examples.
