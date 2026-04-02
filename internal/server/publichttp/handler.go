package publichttp

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"hostbin/internal/domain/documents"
	"hostbin/internal/server/documentsvc"
	"hostbin/internal/server/requestmeta"
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
