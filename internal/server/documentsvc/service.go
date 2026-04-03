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

package documentsvc

import (
	"context"

	"github.com/codelif/hostbin/internal/clock"
	"github.com/codelif/hostbin/internal/domain/documents"
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

func (s *Service) CreateDocument(ctx context.Context, slug string, content []byte) (*documents.Document, error) {
	return s.store.CreateDocument(ctx, slug, content, s.clock.Now())
}

func (s *Service) ReplaceDocument(ctx context.Context, slug string, content []byte) (*documents.Document, error) {
	return s.store.ReplaceDocument(ctx, slug, content, s.clock.Now())
}

func (s *Service) DeleteDocument(ctx context.Context, slug string) error {
	return s.store.DeleteDocument(ctx, slug)
}
