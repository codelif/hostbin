# hostbin

Host-routed plaintext document server with a Gin-based admin API and SQLite storage.

The runnable Go module lives in `server/`.

Quick starts:

- canonical install: `make build`, then `sudo make install-all`
- local development: `make test` or `make run`
- container workflow: `docker compose up --build`
