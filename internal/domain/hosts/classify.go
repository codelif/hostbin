package hosts

import (
	"errors"
	"strconv"
	"strings"

	"hostbin/internal/domain/slugs"
)

var ErrInvalidHost = errors.New("invalid host")

type Kind string

const (
	KindInvalid Kind = "invalid"
	KindAdmin   Kind = "admin"
	KindPublic  Kind = "public"
)

type Info struct {
	Kind Kind
	Host string
	Slug string
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

func ClassifyHost(rawHost, baseDomain, adminHost string, reserved map[string]struct{}) (Info, error) {
	host, err := NormalizeHost(rawHost)
	if err != nil {
		return Info{Kind: KindInvalid}, err
	}

	if host == adminHost {
		return Info{Kind: KindAdmin, Host: host}, nil
	}

	suffix := "." + baseDomain
	if !strings.HasSuffix(host, suffix) {
		return Info{Kind: KindInvalid, Host: host}, ErrInvalidHost
	}

	prefix := strings.TrimSuffix(host, suffix)
	if prefix == "" || strings.Contains(prefix, ".") {
		return Info{Kind: KindInvalid, Host: host}, ErrInvalidHost
	}

	if err := slugs.Validate(prefix, reserved); err != nil {
		return Info{Kind: KindInvalid, Host: host}, ErrInvalidHost
	}

	return Info{Kind: KindPublic, Host: host, Slug: prefix}, nil
}
