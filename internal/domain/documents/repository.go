package documents

import (
	"context"
	"time"
)

type Repository interface {
	ListDocuments(ctx context.Context) ([]DocumentMeta, error)
	GetDocument(ctx context.Context, slug string) (*Document, error)
	GetDocumentMeta(ctx context.Context, slug string) (*DocumentMeta, error)
	CreateDocument(ctx context.Context, slug string, content []byte, now time.Time) (*Document, error)
	ReplaceDocument(ctx context.Context, slug string, content []byte, now time.Time) (*Document, error)
	DeleteDocument(ctx context.Context, slug string) error
}
