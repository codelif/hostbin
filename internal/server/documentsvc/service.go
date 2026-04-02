package documentsvc

import (
	"context"

	"hostbin/internal/clock"
	"hostbin/internal/domain/documents"
)

type Service struct {
	store documents.Repository
	clock clock.Clock
}

func New(store documents.Repository, clock clock.Clock) *Service {
	return &Service{store: store, clock: clock}
}

func (s *Service) ListDocuments(ctx context.Context) ([]documents.DocumentMeta, error) {
	return s.store.ListDocuments(ctx)
}

func (s *Service) GetDocumentMeta(ctx context.Context, slug string) (*documents.DocumentMeta, error) {
	return s.store.GetDocumentMeta(ctx, slug)
}

func (s *Service) GetDocument(ctx context.Context, slug string) (*documents.Document, error) {
	return s.store.GetDocument(ctx, slug)
}

func (s *Service) PutDocument(ctx context.Context, slug string, content []byte) (*documents.Document, error) {
	return s.store.PutDocument(ctx, slug, content, s.clock.Now())
}

func (s *Service) DeleteDocument(ctx context.Context, slug string) error {
	return s.store.DeleteDocument(ctx, slug)
}
