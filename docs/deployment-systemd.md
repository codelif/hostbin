# systemd

Recommended install layout:

- binary: `/usr/local/bin/hostbin`
- env file: `/etc/hostbin/hostbin.env`
- state directory: `/var/lib/hostbin`
- database: `/var/lib/hostbin/data.db`

Setup outline:

1. create a dedicated service account: `sudo useradd --system --home /var/lib/hostbin --shell /usr/sbin/nologin hostbin`
2. build the binary with `make build`
3. install files with `sudo make install-all`
4. edit `/etc/hostbin/hostbin.env` and fill in the secret and domain settings
5. ensure `/etc/hostbin` and `/var/lib/hostbin` are owned by the `hostbin` user
6. enable and start the service with `sudo systemctl enable --now hostbin`

Notes on `make install-all`:

- installs the binary to `/usr/local/bin/hostbin`
- installs the CLI to `/usr/local/bin/hbcli`
- installs the unit to `/etc/systemd/system/hostbin.service`
- installs `/etc/hostbin/hostbin.env` only if it does not already exist
- creates `/var/lib/hostbin`

Verification:

- `sudo systemd-analyze verify /etc/systemd/system/hostbin.service`
- `sudo systemctl status hostbin`
- `curl -H 'Host: admin.domain.com' http://127.0.0.1:8080/api/v1/health`

Operational notes:

- keep `PRESHARED_KEY` only in the env file, not in the unit file
- use filesystem permissions to protect `/etc/hostbin/hostbin.env` and `/var/lib/hostbin/data.db`
- pair this with a reverse proxy that preserves the original `Host` header
- if you raise `MAX_DOC_SIZE`, also raise the reverse proxy request body limit
