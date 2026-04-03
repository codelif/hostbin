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

package dispatch

import (
	"net/http"

	"github.com/codelif/hostbin/internal/domain/hosts"
	"github.com/codelif/hostbin/internal/server/requestmeta"
)

type Dispatcher struct {
	baseDomain string
	adminHost  string
	reserved   map[string]struct{}
	admin      http.Handler
	public     http.Handler
}

func NewHandler(baseDomain, adminHost string, reserved map[string]struct{}, admin, public http.Handler) *Dispatcher {
	return &Dispatcher{
		baseDomain: baseDomain,
		adminHost:  adminHost,
		reserved:   reserved,
		admin:      admin,
		public:     public,
	}
}

func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	info, err := hosts.ClassifyHost(r.Host, d.baseDomain, d.adminHost, d.reserved)
	if meta := requestmeta.FromContext(r.Context()); meta != nil {
		meta.Host = info.Host
		meta.HostKind = info.Kind
		meta.Slug = info.Slug
	}

	if err != nil {
		writePlaintext(w, http.StatusNotFound, "not found\n")
		return
	}

	switch info.Kind {
	case hosts.KindAdmin:
		d.admin.ServeHTTP(w, r)
	case hosts.KindPublic:
		d.public.ServeHTTP(w, r)
	default:
		writePlaintext(w, http.StatusNotFound, "not found\n")
	}
}

func writePlaintext(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}
