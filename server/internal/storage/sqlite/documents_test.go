package sqlite

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"hostbin/internal/storage"
)

func TestDocumentStoreLifecycle(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "data.db"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer db.Close()

	store := NewDocumentStore(db)
	ctx := context.Background()
	createdAt := time.Date(2026, 4, 2, 12, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(10 * time.Minute)

	doc, err := store.PutDocument(ctx, "doc1", []byte("hello"), createdAt)
	if err != nil {
		t.Fatalf("PutDocument(create) error = %v", err)
	}
	if doc.SizeBytes != 5 {
		t.Fatalf("SizeBytes = %d, want 5", doc.SizeBytes)
	}
	if !doc.CreatedAt.Equal(createdAt) {
		t.Fatalf("CreatedAt = %s, want %s", doc.CreatedAt, createdAt)
	}

	doc, err = store.PutDocument(ctx, "doc1", []byte("hello again"), updatedAt)
	if err != nil {
		t.Fatalf("PutDocument(update) error = %v", err)
	}
	if !doc.CreatedAt.Equal(createdAt) {
		t.Fatalf("CreatedAt after update = %s, want %s", doc.CreatedAt, createdAt)
	}
	if !doc.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("UpdatedAt = %s, want %s", doc.UpdatedAt, updatedAt)
	}

	meta, err := store.GetDocumentMeta(ctx, "doc1")
	if err != nil {
		t.Fatalf("GetDocumentMeta() error = %v", err)
	}
	if meta.Slug != "doc1" {
		t.Fatalf("Slug = %q, want doc1", meta.Slug)
	}

	list, err := store.ListDocuments(ctx)
	if err != nil {
		t.Fatalf("ListDocuments() error = %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("len(ListDocuments()) = %d, want 1", len(list))
	}

	if err := store.DeleteDocument(ctx, "doc1"); err != nil {
		t.Fatalf("DeleteDocument() error = %v", err)
	}

	if _, err := store.GetDocument(ctx, "doc1"); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("GetDocument() error = %v, want ErrNotFound", err)
	}
	if err := store.DeleteDocument(ctx, "doc1"); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("DeleteDocument() error = %v, want ErrNotFound", err)
	}
}
