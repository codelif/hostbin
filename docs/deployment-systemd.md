# systemd

Recommended install layout:

- binary: `/usr/local/bin/hostbin`
- env file: `/etc/hostbin/hostbin.env`
- state directory: `/var/lib/hostbin`
- database: `/var/lib/hostbin/data.db`

Setup outline:

1. create a dedicated service account: `sudo useradd --system --home /var/lib/hostbin --shell /usr/sbin/nologin hostbin`
2. install the binary to `/usr/local/bin/hostbin`
3. install `deploy/systemd/hostbin.service` to `/etc/systemd/system/hostbin.service`
4. copy `deploy/systemd/hostbin.env.example` to `/etc/hostbin/hostbin.env` and fill in the secret
5. create `/etc/hostbin` and `/var/lib/hostbin` owned by the `hostbin` user
6. enable and start the service with `sudo systemctl enable --now hostbin`

Verification:

- `sudo systemd-analyze verify /etc/systemd/system/hostbin.service`
- `sudo systemctl status hostbin`
- `curl -H 'Host: admin.domain.com' http://127.0.0.1:8080/api/v1/health`

Operational notes:

- keep `PRESHARED_KEY` only in the env file, not in the unit file
- use filesystem permissions to protect `/etc/hostbin/hostbin.env` and `/var/lib/hostbin/data.db`
- pair this with a reverse proxy that preserves the original `Host` header
- if you raise `MAX_DOC_SIZE`, also raise the reverse proxy request body limit
