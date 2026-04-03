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

package hosts

import (
	"errors"
	"strconv"
	"strings"

	"github.com/codelif/hostbin/internal/domain/slugs"
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
