package integration_test

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"

	"hostbin/internal/app"
	"hostbin/internal/auth"
	"hostbin/internal/clock"
	"hostbin/internal/config"
)

func TestPublicAndAdminFlows(t *testing.T) {
	fixedTime := time.Date(2026, 4, 2, 12, 0, 0, 0, time.UTC)
	handler, cleanup := newTestHandler(t, fixedTime)
	defer cleanup()

	putReq := signedAdminRequest(t, http.MethodPut, "/api/v1/documents/doc1", []byte("hello world"), fixedTime, "aaaa1111bbbb2222cccc3333dddd4444")
	putResp := httptest.NewRecorder()
	handler.ServeHTTP(putResp, putReq)
	if putResp.Code != http.StatusOK {
		t.Fatalf("PUT status = %d, want %d, body=%s", putResp.Code, http.StatusOK, putResp.Body.String())
	}

	var putBody struct {
		Slug      string `json:"slug"`
		URL       string `json:"url"`
		SizeBytes int64  `json:"size_bytes"`
		SHA256    string `json:"sha256"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
	decodeJSON(t, putResp.Body.String(), &putBody)
	if putBody.Slug != "doc1" {
		t.Fatalf("slug = %q, want doc1", putBody.Slug)
	}
	if putBody.URL != "https://doc1.domain.com/" {
		t.Fatalf("url = %q, want https://doc1.domain.com/", putBody.URL)
	}

	publicGetReq := httptest.NewRequest(http.MethodGet, "http://doc1.domain.com/", nil)
	publicGetReq.Host = "doc1.domain.com"
	publicGetResp := httptest.NewRecorder()
	handler.ServeHTTP(publicGetResp, publicGetReq)
	if publicGetResp.Code != http.StatusOK {
		t.Fatalf("public GET status = %d, want %d", publicGetResp.Code, http.StatusOK)
	}
	if body := publicGetResp.Body.String(); body != "hello world" {
		t.Fatalf("public GET body = %q, want %q", body, "hello world")
	}
	etag := publicGetResp.Header().Get("ETag")
	if etag == "" {
		t.Fatal("expected ETag header")
	}

	publicHeadReq := httptest.NewRequest(http.MethodHead, "http://doc1.domain.com/", nil)
	publicHeadReq.Host = "doc1.domain.com"
	publicHeadResp := httptest.NewRecorder()
	handler.ServeHTTP(publicHeadResp, publicHeadReq)
	if publicHeadResp.Code != http.StatusOK {
		t.Fatalf("public HEAD status = %d, want %d", publicHeadResp.Code, http.StatusOK)
	}
	if publicHeadResp.Body.Len() != 0 {
		t.Fatalf("public HEAD body len = %d, want 0", publicHeadResp.Body.Len())
	}

	conditionalReq := httptest.NewRequest(http.MethodGet, "http://doc1.domain.com/", nil)
	conditionalReq.Host = "doc1.domain.com"
	conditionalReq.Header.Set("If-None-Match", etag)
	conditionalResp := httptest.NewRecorder()
	handler.ServeHTTP(conditionalResp, conditionalReq)
	if conditionalResp.Code != http.StatusNotModified {
		t.Fatalf("conditional GET status = %d, want %d", conditionalResp.Code, http.StatusNotModified)
	}

	metadataReq := signedAdminRequest(t, http.MethodGet, "/api/v1/documents/doc1", nil, fixedTime, "bbbb1111bbbb2222cccc3333dddd4444")
	metadataResp := httptest.NewRecorder()
	handler.ServeHTTP(metadataResp, metadataReq)
	if metadataResp.Code != http.StatusOK {
		t.Fatalf("metadata GET status = %d, want %d", metadataResp.Code, http.StatusOK)
	}

	contentReq := signedAdminRequest(t, http.MethodGet, "/api/v1/documents/doc1/content", nil, fixedTime, "cccc1111bbbb2222cccc3333dddd4444")
	contentResp := httptest.NewRecorder()
	handler.ServeHTTP(contentResp, contentReq)
	if contentResp.Code != http.StatusOK {
		t.Fatalf("content GET status = %d, want %d", contentResp.Code, http.StatusOK)
	}
	if contentResp.Body.String() != "hello world" {
		t.Fatalf("content body = %q, want %q", contentResp.Body.String(), "hello world")
	}

	deleteReq := signedAdminRequest(t, http.MethodDelete, "/api/v1/documents/doc1", nil, fixedTime, "dddd1111bbbb2222cccc3333dddd4444")
	deleteResp := httptest.NewRecorder()
	handler.ServeHTTP(deleteResp, deleteReq)
	if deleteResp.Code != http.StatusOK {
		t.Fatalf("DELETE status = %d, want %d", deleteResp.Code, http.StatusOK)
	}

	missingReq := httptest.NewRequest(http.MethodGet, "http://doc1.domain.com/", nil)
	missingReq.Host = "doc1.domain.com"
	missingResp := httptest.NewRecorder()
	handler.ServeHTTP(missingResp, missingReq)
	if missingResp.Code != http.StatusNotFound || missingResp.Body.String() != "not found\n" {
		t.Fatalf("missing public response = (%d, %q), want (404, %q)", missingResp.Code, missingResp.Body.String(), "not found\n")
	}

	invalidHostReq := httptest.NewRequest(http.MethodGet, "http://domain.com/", nil)
	invalidHostReq.Host = "domain.com"
	invalidHostResp := httptest.NewRecorder()
	handler.ServeHTTP(invalidHostResp, invalidHostReq)
	if invalidHostResp.Code != http.StatusNotFound || invalidHostResp.Body.String() != "not found\n" {
		t.Fatalf("invalid host response = (%d, %q), want (404, %q)", invalidHostResp.Code, invalidHostResp.Body.String(), "not found\n")
	}

	methodReq := httptest.NewRequest(http.MethodPost, "http://admin.domain.com/api/v1/documents", nil)
	methodReq.Host = "admin.domain.com"
	methodResp := httptest.NewRecorder()
	handler.ServeHTTP(methodResp, methodReq)
	if methodResp.Code != http.StatusMethodNotAllowed || strings.TrimSpace(methodResp.Body.String()) != `{"error":"method_not_allowed"}` {
		t.Fatalf("method response = (%d, %q)", methodResp.Code, methodResp.Body.String())
	}

	utf8Req := signedAdminRequest(t, http.MethodPut, "/api/v1/documents/doc2", []byte{0xff, 0xfe}, fixedTime, "eeee1111bbbb2222cccc3333dddd4444")
	utf8Resp := httptest.NewRecorder()
	handler.ServeHTTP(utf8Resp, utf8Req)
	if utf8Resp.Code != http.StatusBadRequest || strings.TrimSpace(utf8Resp.Body.String()) != `{"error":"invalid_utf8"}` {
		t.Fatalf("invalid utf8 response = (%d, %q)", utf8Resp.Code, utf8Resp.Body.String())
	}
}

func newTestHandler(t *testing.T, fixedTime time.Time) (http.Handler, func()) {
	t.Helper()
	cfg := config.Config{
		ListenAddr:         "127.0.0.1:8080",
		BaseDomain:         "domain.com",
		AdminHost:          "admin.domain.com",
		PresharedKey:       "01234567890123456789012345678901",
		DBPath:             filepath.Join(t.TempDir(), "data.db"),
		ReservedSubdomains: []string{"admin", "www", "api"},
		ReservedSet: map[string]struct{}{
			"admin": {},
			"www":   {},
			"api":   {},
		},
		MaxDocSize:        1024,
		AuthTimestampSkew: time.Minute,
		NonceTTL:          5 * time.Minute,
		TrustedProxyCIDRs: []string{"127.0.0.1/32"},
		LogLevel:          "info",
	}

	application, err := app.New(cfg, app.Options{Clock: clock.Fixed{Time: fixedTime}, Logger: zap.NewNop()})
	if err != nil {
		t.Fatalf("app.New() error = %v", err)
	}

	return application.Server.Handler, func() {
		_ = application.Close()
	}
}

func signedAdminRequest(t *testing.T, method, path string, body []byte, timestamp time.Time, nonceValue string) *http.Request {
	t.Helper()
	var bodyReader *strings.Reader
	if body == nil {
		bodyReader = strings.NewReader("")
	} else {
		bodyReader = strings.NewReader(string(body))
	}

	req := httptest.NewRequest(method, "http://admin.domain.com"+path, bodyReader)
	req.Host = "admin.domain.com"
	req.URL.RawPath = path
	req.Header.Set(auth.HeaderTimestamp, strconv.FormatInt(timestamp.Unix(), 10))
	req.Header.Set(auth.HeaderNonce, nonceValue)
	canonical := auth.CanonicalRequest(req, auth.SHA256Hex(body), req.Header.Get(auth.HeaderTimestamp), nonceValue)
	req.Header.Set(auth.HeaderSignature, hex.EncodeToString(auth.Sign([]byte("01234567890123456789012345678901"), canonical)))
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	return req
}

func decodeJSON(t *testing.T, raw string, target any) {
	t.Helper()
	if err := json.Unmarshal([]byte(raw), target); err != nil {
		t.Fatalf("json.Unmarshal(%q) error = %v", raw, err)
	}
}
