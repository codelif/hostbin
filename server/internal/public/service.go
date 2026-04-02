package public

import (
	"context"

	"hostbin/internal/model"
	"hostbin/internal/storage"
)

type Service struct {
	store storage.DocumentStore
}

func NewService(store storage.DocumentStore) *Service {
	return &Service{store: store}
}

func (s *Service) GetDocument(ctx context.Context, slug string) (*model.Document, error) {
	return s.store.GetDocument(ctx, slug)
}
