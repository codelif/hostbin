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

package publichttp

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/codelif/hostbin/internal/domain/documents"
	"github.com/codelif/hostbin/internal/server/documentsvc"
	"github.com/codelif/hostbin/internal/server/requestmeta"
)

type Handler struct {
	service *documentsvc.Service
}

func NewHandler(service *documentsvc.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetRoot(c *gin.Context) {
	meta := requestmeta.FromContext(c.Request.Context())
	if meta == nil || meta.Slug == "" {
		writePlaintextError(c, http.StatusNotFound, "not found\n")
		return
	}

	doc, err := h.service.GetDocument(c.Request.Context(), meta.Slug)
	if err != nil {
		if errors.Is(err, documents.ErrNotFound) {
			writePlaintextError(c, http.StatusNotFound, "not found\n")
			return
		}

		writePlaintextError(c, http.StatusInternalServerError, "internal error\n")
		return
	}

	etag := fmt.Sprintf(`"sha256-%s"`, doc.SHA256)
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("ETag", etag)

	if c.GetHeader("If-None-Match") == etag {
		c.Status(http.StatusNotModified)
		return
	}

	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusOK)
		return
	}

	c.Data(http.StatusOK, "text/plain; charset=utf-8", doc.Content)
}

func writePlaintextError(c *gin.Context, status int, body string) {
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("X-Content-Type-Options", "nosniff")
	c.String(status, body)
}
