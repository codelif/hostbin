package adminhttp

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"hostbin/internal/domain/documents"
	"hostbin/internal/domain/hosts"
	"hostbin/internal/domain/slugs"
	"hostbin/internal/protocol/adminv1"
	"hostbin/internal/server/documentsvc"
)

type Handler struct {
	service    *documentsvc.Service
	baseDomain string
	reserved   map[string]struct{}
}

func NewHandler(service *documentsvc.Service, baseDomain string, reserved map[string]struct{}) *Handler {
	return &Handler{service: service, baseDomain: baseDomain, reserved: reserved}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, adminv1.HealthResponse{Status: "ok"})
}

func (h *Handler) ListDocuments(c *gin.Context) {
	docs, err := h.service.ListDocuments(c.Request.Context())
	if err != nil {
		writeError(c, http.StatusInternalServerError, adminv1.ErrorInternal)
		return
	}

	items := make([]adminv1.DocumentResponse, 0, len(docs))
	for _, doc := range docs {
		items = append(items, toDocumentResponse(h.baseDomain, doc))
	}

	c.JSON(http.StatusOK, adminv1.ListDocumentsResponse{Documents: items})
}

func (h *Handler) GetDocument(c *gin.Context) {
	slugValue, ok := h.validatedSlug(c)
	if !ok {
		return
	}

	doc, err := h.service.GetDocumentMeta(c.Request.Context(), slugValue)
	if err != nil {
		h.handleStoreError(c, err)
		return
	}

	c.JSON(http.StatusOK, toDocumentResponse(h.baseDomain, *doc))
}

func (h *Handler) GetDocumentContent(c *gin.Context) {
	slugValue, ok := h.validatedSlug(c)
	if !ok {
		return
	}

	doc, err := h.service.GetDocument(c.Request.Context(), slugValue)
	if err != nil {
		h.handleStoreError(c, err)
		return
	}

	c.Data(http.StatusOK, "text/plain; charset=utf-8", doc.Content)
}

func (h *Handler) PutDocument(c *gin.Context) {
	slugValue, ok := h.validatedSlug(c)
	if !ok {
		return
	}

	if !validTextPlainContentType(c.GetHeader("Content-Type")) {
		writeError(c, http.StatusBadRequest, adminv1.ErrorBadRequest)
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		if isBodyTooLarge(err) {
			writeError(c, http.StatusRequestEntityTooLarge, adminv1.ErrorDocumentTooLarge)
			return
		}

		writeError(c, http.StatusBadRequest, adminv1.ErrorBadRequest)
		return
	}

	if !utf8.Valid(body) {
		writeError(c, http.StatusBadRequest, adminv1.ErrorInvalidUTF8)
		return
	}

	doc, err := h.service.PutDocument(c.Request.Context(), slugValue, body)
	if err != nil {
		writeError(c, http.StatusInternalServerError, adminv1.ErrorInternal)
		return
	}

	c.JSON(http.StatusOK, toDocumentResponse(h.baseDomain, toMeta(*doc)))
}

func (h *Handler) DeleteDocument(c *gin.Context) {
	slugValue, ok := h.validatedSlug(c)
	if !ok {
		return
	}

	err := h.service.DeleteDocument(c.Request.Context(), slugValue)
	if err != nil {
		h.handleStoreError(c, err)
		return
	}

	c.JSON(http.StatusOK, adminv1.DeleteResponse{Deleted: true, Slug: slugValue})
}

func (h *Handler) validatedSlug(c *gin.Context) (string, bool) {
	slugValue := strings.ToLower(strings.TrimSpace(c.Param("slug")))
	if err := slugs.Validate(slugValue, h.reserved); err != nil {
		writeError(c, http.StatusBadRequest, adminv1.ErrorInvalidSlug)
		return "", false
	}

	return slugValue, true
}

func (h *Handler) handleStoreError(c *gin.Context, err error) {
	if errors.Is(err, documents.ErrNotFound) {
		writeError(c, http.StatusNotFound, adminv1.ErrorNotFound)
		return
	}

	writeError(c, http.StatusInternalServerError, adminv1.ErrorInternal)
}

func validTextPlainContentType(raw string) bool {
	mediaType, params, err := mime.ParseMediaType(raw)
	if err != nil || mediaType != "text/plain" {
		return false
	}

	if charset, ok := params["charset"]; ok && !strings.EqualFold(charset, "utf-8") {
		return false
	}

	return true
}

func isBodyTooLarge(err error) bool {
	var target *http.MaxBytesError
	return errors.As(err, &target)
}

func toMeta(doc documents.Document) documents.DocumentMeta {
	return documents.DocumentMeta{
		Slug:      doc.Slug,
		SHA256:    doc.SHA256,
		SizeBytes: doc.SizeBytes,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}

func toDocumentResponse(baseDomain string, doc documents.DocumentMeta) adminv1.DocumentResponse {
	return adminv1.DocumentResponse{
		Slug:      doc.Slug,
		URL:       hosts.DocumentURL(baseDomain, doc.Slug),
		SizeBytes: doc.SizeBytes,
		SHA256:    doc.SHA256,
		CreatedAt: doc.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: doc.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func writeError(c *gin.Context, status int, code string) {
	c.AbortWithStatusJSON(status, adminv1.ErrorResponse{Error: code})
}
