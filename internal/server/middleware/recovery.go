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

package middleware

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/codelif/hostbin/internal/domain/hosts"
	"github.com/codelif/hostbin/internal/server/requestmeta"
)

func Recovery(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				meta := requestmeta.FromContext(r.Context())
				logger.Error("panic recovered",
					zap.Any("panic", recovered),
					zap.String("request_id", metaValue(meta, func(m *requestmeta.Meta) string { return m.RequestID })),
					zap.String("host", metaValue(meta, func(m *requestmeta.Meta) string { return m.Host })),
					zap.String("path", r.URL.Path),
				)

				if meta != nil && meta.HostKind == hosts.KindAdmin {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"error":"internal_error"}`))
					return
				}

				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("internal error\n"))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func metaValue(meta *requestmeta.Meta, extract func(*requestmeta.Meta) string) string {
	if meta == nil {
		return ""
	}
	return extract(meta)
}
