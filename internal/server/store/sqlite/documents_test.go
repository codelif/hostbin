// Copyright (c) 2026 Harsh Sharma <harsh@codelif.in>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// SPDX-License-Identifier: MIT

package sqlite

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/codelif/hostbin/internal/domain/documents"
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

	doc, err := store.CreateDocument(ctx, "doc1", []byte("hello"), createdAt)
	if err != nil {
		t.Fatalf("CreateDocument() error = %v", err)
	}
	if doc.SizeBytes != 5 {
		t.Fatalf("SizeBytes = %d, want 5", doc.SizeBytes)
	}
	if !doc.CreatedAt.Equal(createdAt) {
		t.Fatalf("CreatedAt = %s, want %s", doc.CreatedAt, createdAt)
	}

	if _, err := store.CreateDocument(ctx, "doc1", []byte("hello again"), updatedAt); !errors.Is(err, documents.ErrAlreadyExists) {
		t.Fatalf("CreateDocument(duplicate) error = %v, want ErrAlreadyExists", err)
	}

	doc, err = store.ReplaceDocument(ctx, "doc1", []byte("hello again"), updatedAt)
	if err != nil {
		t.Fatalf("ReplaceDocument() error = %v", err)
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

	if _, err := store.ReplaceDocument(ctx, "missing", []byte("hello"), updatedAt); !errors.Is(err, documents.ErrNotFound) {
		t.Fatalf("ReplaceDocument(missing) error = %v, want ErrNotFound", err)
	}

	if err := store.DeleteDocument(ctx, "doc1"); err != nil {
		t.Fatalf("DeleteDocument() error = %v", err)
	}

	if _, err := store.GetDocument(ctx, "doc1"); !errors.Is(err, documents.ErrNotFound) {
		t.Fatalf("GetDocument() error = %v, want ErrNotFound", err)
	}
	if err := store.DeleteDocument(ctx, "doc1"); !errors.Is(err, documents.ErrNotFound) {
		t.Fatalf("DeleteDocument() error = %v, want ErrNotFound", err)
	}
}
