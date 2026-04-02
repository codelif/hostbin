package storage

import (
	"context"
	"time"

	"hostbin/internal/model"
)

type DocumentStore interface {
	ListDocuments(ctx context.Context) ([]model.DocumentMeta, error)
	GetDocument(ctx context.Context, slug string) (*model.Document, error)
	GetDocumentMeta(ctx context.Context, slug string) (*model.DocumentMeta, error)
	PutDocument(ctx context.Context, slug string, content []byte, now time.Time) (*model.Document, error)
	DeleteDocument(ctx context.Context, slug string) error
}
