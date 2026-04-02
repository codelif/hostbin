package requestmeta

import (
	"context"

	"hostbin/internal/domain/hosts"
)

type Meta struct {
	RequestID string
	Host      string
	HostKind  hosts.Kind
	Slug      string
}

type contextKey struct{}

func WithContext(ctx context.Context, meta *Meta) context.Context {
	return context.WithValue(ctx, contextKey{}, meta)
}

func FromContext(ctx context.Context) *Meta {
	meta, _ := ctx.Value(contextKey{}).(*Meta)
	return meta
}
