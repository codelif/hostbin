package admin

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"hostbin/internal/model"
	"hostbin/internal/slug"
	"hostbin/internal/storage"
)

const (
	errorUnauthorized     = "unauthorized"
	errorInvalidSignature = "invalid_signature"
	errorInvalidTimestamp = "invalid_timestamp"
	errorReplayedNonce    = "replayed_nonce"
	errorInvalidSlug      = "invalid_slug"
	errorNotFound         = "not_found"
	errorDocTooLarge      = "document_too_large"
	errorBadRequest       = "bad_request"
	errorInvalidUTF8      = "invalid_utf8"
	errorMethodNotAllowed = "method_not_allowed"
	errorInternal         = "internal_error"
)

type Handler struct {
	service  *Service
	reserved map[string]struct{}
}

func NewHandler(service *Service, reserved map[string]struct{}) *Handler {
	return &Handler{service: service, reserved: reserved}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}

func (h *Handler) ListDocuments(c *gin.Context) {
	docs, err := h.service.ListDocuments(c.Request.Context())
	if err != nil {
		writeError(c, http.StatusInternalServerError, errorInternal)
		return
	}

	items := make([]DocumentResponse, 0, len(docs))
	for _, doc := range docs {
		items = append(items, toDocumentResponse(h.service.baseDomain, doc))
	}

	c.JSON(http.StatusOK, ListDocumentsResponse{Documents: items})
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

	c.JSON(http.StatusOK, toDocumentResponse(h.service.baseDomain, *doc))
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
		writeError(c, http.StatusBadRequest, errorBadRequest)
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		if isBodyTooLarge(err) {
			writeError(c, http.StatusRequestEntityTooLarge, errorDocTooLarge)
			return
		}

		writeError(c, http.StatusBadRequest, errorBadRequest)
		return
	}

	if !utf8.Valid(body) {
		writeError(c, http.StatusBadRequest, errorInvalidUTF8)
		return
	}

	doc, err := h.service.PutDocument(c.Request.Context(), slugValue, body)
	if err != nil {
		writeError(c, http.StatusInternalServerError, errorInternal)
		return
	}

	c.JSON(http.StatusOK, toDocumentResponse(h.service.baseDomain, toMeta(*doc)))
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

	c.JSON(http.StatusOK, DeleteResponse{Deleted: true, Slug: slugValue})
}

func (h *Handler) validatedSlug(c *gin.Context) (string, bool) {
	slugValue := strings.ToLower(strings.TrimSpace(c.Param("slug")))
	if err := slug.Validate(slugValue, h.reserved); err != nil {
		writeError(c, http.StatusBadRequest, errorInvalidSlug)
		return "", false
	}

	return slugValue, true
}

func (h *Handler) handleStoreError(c *gin.Context, err error) {
	if errors.Is(err, storage.ErrNotFound) {
		writeError(c, http.StatusNotFound, errorNotFound)
		return
	}

	writeError(c, http.StatusInternalServerError, errorInternal)
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

func toMeta(doc model.Document) model.DocumentMeta {
	return model.DocumentMeta{
		Slug:      doc.Slug,
		SHA256:    doc.SHA256,
		SizeBytes: doc.SizeBytes,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}

func writeError(c *gin.Context, status int, code string) {
	c.AbortWithStatusJSON(status, ErrorResponse{Error: code})
}
