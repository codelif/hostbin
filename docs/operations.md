# Operations

This page covers day-2 operational tasks for a running `hostbin` deployment.

Use it after the service is already live and you need to restart, back up, restore, upgrade, or verify the deployment.

## Logs

For `systemd` deployments, follow the application logs with:

```bash
sudo journalctl -u hostbin -f
```

If you run Caddy as a separate service:

```bash
sudo journalctl -u caddy -f
```

## Service control

Restart the app after changing `/etc/hostbin/hostbin.env`:

```bash
sudo systemctl restart hostbin
```

Restart Caddy after changing its config:

```bash
sudo systemctl restart caddy
```

Inspect status:

```bash
sudo systemctl status hostbin --no-pager
sudo systemctl status caddy --no-pager
```

## Database location

Typical production database path:

```text
/var/lib/hostbin/data.db
```

Your actual location is controlled by `DB_PATH`.

## Backups

The safest simple backup flow is stop, copy, start:

```bash
sudo systemctl stop hostbin
sudo cp /var/lib/hostbin/data.db /var/lib/hostbin/data.db.bak.$(date +%F-%H%M%S)
sudo systemctl start hostbin
```

Store backups outside the VM as well if the deployment matters.

## Restore

Restore by replacing the database file while the service is stopped:

```bash
sudo systemctl stop hostbin
sudo cp /path/to/backup.db /var/lib/hostbin/data.db
sudo chown hostbin:hostbin /var/lib/hostbin/data.db
sudo chmod 600 /var/lib/hostbin/data.db
sudo systemctl start hostbin
```

After restore, verify with:

- `curl https://hbadmin.example.com/api/v1/health`
- `hbcli list`

## Upgrades

From the repository root:

```bash
git pull
make build
sudo make install-all
sudo systemctl restart hostbin
```

If you also rebuild Caddy or change its config, restart Caddy separately.

## Secret rotation

To rotate `PRESHARED_KEY`:

1. update `/etc/hostbin/hostbin.env`
2. restart `hostbin`
3. update every `hbcli` config that uses the old key
4. run `hbcli config check`

During rotation, old clients fail auth until updated.

## Capacity and limits

- `MAX_DOC_SIZE` controls the app-side document size limit
- your reverse proxy body limit must match or exceed it
- SQLite is simple and effective, but storage, backup, and I/O sizing still matter on very busy hosts

## Host and proxy changes

If you change any of these values:

- `BASE_DOMAIN`
- `ADMIN_HOST`
- `RESERVED_SUBDOMAINS`
- `TRUST_PROXY_HEADERS`
- `TRUSTED_PROXY_CIDRS`

restart the service and re-run your smoke checks.

## Suggested smoke checks after any change

Admin health:

```bash
curl -i https://hbadmin.example.com/api/v1/health
```

CLI auth check:

```bash
hbcli config check
```

Public fetch:

```bash
curl -i https://hello.example.com/
```

## Related docs

- deployment index: `docs/deployment.md`
- production deployment: `docs/deployment-cloudflare-caddy-systemd.md`
- troubleshooting: `docs/troubleshooting.md`
