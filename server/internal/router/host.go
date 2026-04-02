package router

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"hostbin/internal/slug"
)

var ErrInvalidHost = errors.New("invalid host")

type HostKind string

const (
	HostKindInvalid HostKind = "invalid"
	HostKindAdmin   HostKind = "admin"
	HostKindPublic  HostKind = "public"
)

type HostInfo struct {
	Kind HostKind
	Host string
	Slug string
}

type RequestMeta struct {
	RequestID string
	Host      string
	HostKind  string
	Slug      string
}

type requestMetaKey struct{}

func WithRequestMeta(ctx context.Context, meta *RequestMeta) context.Context {
	return context.WithValue(ctx, requestMetaKey{}, meta)
}

func RequestMetaFromContext(ctx context.Context) *RequestMeta {
	meta, _ := ctx.Value(requestMetaKey{}).(*RequestMeta)
	return meta
}

func NormalizeHost(rawHost string) (string, error) {
	host := strings.TrimSpace(strings.ToLower(rawHost))
	if host == "" {
		return "", ErrInvalidHost
	}

	if i := strings.LastIndex(host, ":"); i > 0 && !strings.Contains(host[i+1:], ":") {
		if _, err := strconv.Atoi(host[i+1:]); err == nil {
			host = host[:i]
		}
	}

	host = strings.TrimSuffix(host, ".")
	if host == "" || strings.Contains(host, "[") || strings.Contains(host, "]") {
		return "", ErrInvalidHost
	}

	return host, nil
}

func ClassifyHost(rawHost, baseDomain, adminHost string, reserved map[string]struct{}) (HostInfo, error) {
	host, err := NormalizeHost(rawHost)
	if err != nil {
		return HostInfo{Kind: HostKindInvalid}, err
	}

	if host == adminHost {
		return HostInfo{Kind: HostKindAdmin, Host: host}, nil
	}

	suffix := "." + baseDomain
	if !strings.HasSuffix(host, suffix) {
		return HostInfo{Kind: HostKindInvalid, Host: host}, ErrInvalidHost
	}

	prefix := strings.TrimSuffix(host, suffix)
	if prefix == "" || strings.Contains(prefix, ".") {
		return HostInfo{Kind: HostKindInvalid, Host: host}, ErrInvalidHost
	}

	if err := slug.Validate(prefix, reserved); err != nil {
		return HostInfo{Kind: HostKindInvalid, Host: host}, ErrInvalidHost
	}

	return HostInfo{Kind: HostKindPublic, Host: host, Slug: prefix}, nil
}
