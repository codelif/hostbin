package sqlite

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"hostbin/internal/domain/documents"
)

type DocumentStore struct {
	db *sql.DB
}

func NewDocumentStore(db *sql.DB) *DocumentStore {
	return &DocumentStore{db: db}
}

func (s *DocumentStore) ListDocuments(ctx context.Context) ([]documents.DocumentMeta, error) {
	const query = `
SELECT slug, content_sha256, size_bytes, created_at, updated_at
FROM documents
ORDER BY updated_at DESC, slug ASC`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []documents.DocumentMeta
	for rows.Next() {
		var (
			doc          documents.DocumentMeta
			createdAtRaw string
			updatedAtRaw string
		)

		if err := rows.Scan(&doc.Slug, &doc.SHA256, &doc.SizeBytes, &createdAtRaw, &updatedAtRaw); err != nil {
			return nil, err
		}

		doc.CreatedAt, err = parseTimestamp(createdAtRaw)
		if err != nil {
			return nil, err
		}

		doc.UpdatedAt, err = parseTimestamp(updatedAtRaw)
		if err != nil {
			return nil, err
		}

		docs = append(docs, doc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return docs, nil
}

func (s *DocumentStore) GetDocument(ctx context.Context, slug string) (*documents.Document, error) {
	const query = `
SELECT slug, content, content_sha256, size_bytes, created_at, updated_at
FROM documents
WHERE slug = ?`

	var (
		doc          documents.Document
		createdAtRaw string
		updatedAtRaw string
	)

	err := s.db.QueryRowContext(ctx, query, slug).Scan(
		&doc.Slug,
		&doc.Content,
		&doc.SHA256,
		&doc.SizeBytes,
		&createdAtRaw,
		&updatedAtRaw,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, documents.ErrNotFound
		}
		return nil, err
	}

	doc.CreatedAt, err = parseTimestamp(createdAtRaw)
	if err != nil {
		return nil, err
	}

	doc.UpdatedAt, err = parseTimestamp(updatedAtRaw)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (s *DocumentStore) GetDocumentMeta(ctx context.Context, slug string) (*documents.DocumentMeta, error) {
	const query = `
SELECT slug, content_sha256, size_bytes, created_at, updated_at
FROM documents
WHERE slug = ?`

	var (
		doc          documents.DocumentMeta
		createdAtRaw string
		updatedAtRaw string
	)

	err := s.db.QueryRowContext(ctx, query, slug).Scan(
		&doc.Slug,
		&doc.SHA256,
		&doc.SizeBytes,
		&createdAtRaw,
		&updatedAtRaw,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, documents.ErrNotFound
		}
		return nil, err
	}

	doc.CreatedAt, err = parseTimestamp(createdAtRaw)
	if err != nil {
		return nil, err
	}

	doc.UpdatedAt, err = parseTimestamp(updatedAtRaw)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (s *DocumentStore) CreateDocument(ctx context.Context, slug string, content []byte, now time.Time) (*documents.Document, error) {
	now = now.UTC().Truncate(time.Second)
	nowRaw := now.Format(time.RFC3339)
	hashHex, sizeBytes := contentMetadata(content)

	const query = `
INSERT INTO documents (slug, content, content_sha256, size_bytes, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)`

	if _, err := s.db.ExecContext(ctx, query, slug, content, hashHex, sizeBytes, nowRaw, nowRaw); err != nil {
		if isUniqueConstraint(err) {
			return nil, documents.ErrAlreadyExists
		}
		return nil, err
	}

	return s.GetDocument(ctx, slug)
}

func (s *DocumentStore) ReplaceDocument(ctx context.Context, slug string, content []byte, now time.Time) (*documents.Document, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	now = now.UTC().Truncate(time.Second)
	nowRaw := now.Format(time.RFC3339)
	hashHex, sizeBytes := contentMetadata(content)

	const existingQuery = `SELECT created_at FROM documents WHERE slug = ?`
	var createdAtRaw string
	err = tx.QueryRowContext(ctx, existingQuery, slug).Scan(&createdAtRaw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, documents.ErrNotFound
		}
		return nil, err
	}

	const update = `
UPDATE documents
SET content = ?, content_sha256 = ?, size_bytes = ?, updated_at = ?
WHERE slug = ?`

	result, err := tx.ExecContext(ctx, update, content, hashHex, sizeBytes, nowRaw, slug)
	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, documents.ErrNotFound
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetDocument(ctx, slug)
}

func (s *DocumentStore) DeleteDocument(ctx context.Context, slug string) error {
	const query = `DELETE FROM documents WHERE slug = ?`

	result, err := s.db.ExecContext(ctx, query, slug)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return documents.ErrNotFound
	}

	return nil
}

func parseTimestamp(raw string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse timestamp %q: %w", raw, err)
	}

	return parsed.UTC(), nil
}

func contentMetadata(content []byte) (string, int64) {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:]), int64(len(content))
}

func isUniqueConstraint(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "unique")
}
