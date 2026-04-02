package admin

import (
	"context"
	"time"

	"hostbin/internal/clock"
	"hostbin/internal/model"
	"hostbin/internal/router"
	"hostbin/internal/storage"
)

type Service struct {
	store      storage.DocumentStore
	baseDomain string
	clock      clock.Clock
}

func NewService(store storage.DocumentStore, baseDomain string, clock clock.Clock) *Service {
	return &Service{
		store:      store,
		baseDomain: baseDomain,
		clock:      clock,
	}
}

func (s *Service) ListDocuments(ctx context.Context) ([]model.DocumentMeta, error) {
	return s.store.ListDocuments(ctx)
}

func (s *Service) GetDocumentMeta(ctx context.Context, slug string) (*model.DocumentMeta, error) {
	return s.store.GetDocumentMeta(ctx, slug)
}

func (s *Service) GetDocument(ctx context.Context, slug string) (*model.Document, error) {
	return s.store.GetDocument(ctx, slug)
}

func (s *Service) PutDocument(ctx context.Context, slug string, content []byte) (*model.Document, error) {
	return s.store.PutDocument(ctx, slug, content, s.clock.Now())
}

func (s *Service) DeleteDocument(ctx context.Context, slug string) error {
	return s.store.DeleteDocument(ctx, slug)
}

func (s *Service) DocumentURL(slug string) string {
	return router.DocumentURL(s.baseDomain, slug)
}

func toDocumentResponse(baseDomain string, doc model.DocumentMeta) DocumentResponse {
	return DocumentResponse{
		Slug:      doc.Slug,
		URL:       router.DocumentURL(baseDomain, doc.Slug),
		SizeBytes: doc.SizeBytes,
		SHA256:    doc.SHA256,
		CreatedAt: doc.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: doc.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
