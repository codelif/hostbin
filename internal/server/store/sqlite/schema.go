package sqlite

import (
	"context"
	"database/sql"
)

const schemaDocuments = `
CREATE TABLE IF NOT EXISTS documents (
    slug TEXT PRIMARY KEY,
    content BLOB NOT NULL,
    content_sha256 TEXT NOT NULL,
    size_bytes INTEGER NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_documents_updated_at
ON documents(updated_at DESC);
`

func initSchema(ctx context.Context, db *sql.DB) error {
	statements := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA busy_timeout = 5000;",
		"PRAGMA synchronous = NORMAL;",
		schemaDocuments,
	}

	for _, stmt := range statements {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}

	return nil
}
