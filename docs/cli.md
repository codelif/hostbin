# CLI

This page documents `hbcli`, the bundled client for configuring access to a `hostbin` server and managing documents through the admin API.

## What `hbcli` does

`hbcli` is the easiest way to work with the admin API because it:

- stores your admin server URL and preshared key locally
- signs authenticated requests automatically
- supports interactive and non-interactive document workflows
- prints document metadata and content in operator-friendly formats

Because the server routes by hostname, `hbcli` works best when your local development setup uses a wildcard loopback domain such as `lvh.me` instead of plain `localhost`.

## Install

Build locally:

```bash
make build
```

Then use `./bin/hbcli`, or install system-wide:

```bash
sudo make install-all
```

This installs `hbcli` to `/usr/local/bin/hbcli`.

## Configuration file

Default config path:

```text
~/.config/hbcli/config.toml
```

Override it for any command with:

```bash
hbcli --config /path/to/config.toml <command>
```

### Config fields

- `server_url` - admin base URL such as `https://hbadmin.example.com`
- `auth_key` - preshared secret used to sign admin requests
- `timeout` - HTTP timeout, default `10s`
- `editor` - editor command for interactive edit flows
- `color` - `auto`, `always`, or `never`

Example file:

```toml
server_url = "https://hbadmin.example.com"
auth_key = "replace-with-a-long-random-secret"
timeout = "10s"
editor = "vim"
color = "auto"
```

## Bootstrap config

Interactive:

```bash
hbcli config init
```

Non-interactive:

```bash
hbcli config init \
  --server-url https://hbadmin.example.com \
  --auth-key 'your-preshared-key'
```

Read the key from stdin instead of the shell history:

```bash
printf '%s' "$PRESHARED_KEY" | hbcli config init \
  --server-url https://hbadmin.example.com \
  --auth-key-stdin
```

Verify the saved config:

```bash
hbcli config check
```

This checks both:

- server reachability via `/api/v1/health`
- signing/authentication via `/api/v1/auth/check`

Local development example:

```bash
hbcli config init \
  --server-url http://hbadmin.lvh.me:8080 \
  --auth-key 'your-preshared-key'
```

## Daily commands

List documents:

```bash
hbcli list
```

Show one document's metadata:

```bash
hbcli info hello
```

Print raw content to stdout:

```bash
hbcli get hello
```

Save raw content to a file:

```bash
hbcli get hello --save ./hello.txt
```

Create a new document from stdin:

```bash
printf 'hello from hostbin\n' | hbcli new hello
```

Create a new document from a file:

```bash
hbcli new hello ./hello.txt
```

Replace an existing document from stdin:

```bash
printf 'updated content\n' | hbcli put hello
```

Edit an existing document in your editor:

```bash
hbcli edit hello
```

Delete a document interactively:

```bash
hbcli delete hello
```

Delete a document non-interactively:

```bash
hbcli delete hello --yes
```

## Interactive behavior

- `hbcli new <slug>` with no file and an interactive terminal opens your editor
- `hbcli put <slug>` with no file and an interactive terminal opens your editor
- `hbcli edit <slug>` fetches the current content, opens it in your editor, and uploads changes only if the content changed
- `hbcli delete <slug>` prompts unless you pass `--yes`

## Non-interactive behavior

- `hbcli new <slug>` reads from stdin when stdin is not a TTY
- `hbcli put <slug>` reads from stdin when stdin is not a TTY
- `hbcli delete <slug>` requires `--yes` when stdin is not interactive

These behaviors make the CLI easy to use in shell pipelines and automation.

## Common flows

Create, inspect, and fetch a document:

```bash
printf 'hello\n' | hbcli new hello
hbcli info hello
hbcli get hello
```

Rotate to a new config file for another environment:

```bash
hbcli --config ~/.config/hbcli/staging.toml config init \
  --server-url https://hbadmin.staging.example.com \
  --auth-key 'staging-secret'
```

## Troubleshooting

- `no configuration found` -> run `hbcli config init`
- `health check failed` -> confirm `server_url` points to the exact admin host and that DNS or `/etc/hosts` resolves it
- `auth check failed` -> confirm `auth_key` matches the server `PRESHARED_KEY`
- `invalid slug` -> use lowercase slug labels only; reserved labels are also rejected
- localhost-only setup feels broken -> use `lvh.me` or another wildcard loopback domain so `hbcli` can reach the real admin hostname

More operational troubleshooting lives in `docs/troubleshooting.md`.
