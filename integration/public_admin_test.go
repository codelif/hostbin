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

	"github.com/codelif/hostbin/internal/clock"
	"github.com/codelif/hostbin/internal/protocol/authsig"
	"github.com/codelif/hostbin/internal/server/app"
	"github.com/codelif/hostbin/internal/server/config"
)

func TestPublicAndAdminFlows(t *testing.T) {
	fixedTime := time.Date(2026, 4, 2, 12, 0, 0, 0, time.UTC)
	handler, cleanup := newTestHandler(t, fixedTime)
	defer cleanup()

	createReq := signedAdminRequest(t, http.MethodPost, "/api/v1/documents/doc1", []byte("hello world"), fixedTime, "aaaa1111bbbb2222cccc3333dddd4444")
	createResp := httptest.NewRecorder()
	handler.ServeHTTP(createResp, createReq)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, want %d, body=%s", createResp.Code, http.StatusCreated, createResp.Body.String())
	}
	if cacheControl := createResp.Header().Get("Cache-Control"); cacheControl != "no-store" {
		t.Fatalf("POST Cache-Control = %q, want %q", cacheControl, "no-store")
	}

	var putBody struct {
		Slug      string `json:"slug"`
		URL       string `json:"url"`
		SizeBytes int64  `json:"size_bytes"`
		SHA256    string `json:"sha256"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
	decodeJSON(t, createResp.Body.String(), &putBody)
	if putBody.Slug != "doc1" {
		t.Fatalf("slug = %q, want doc1", putBody.Slug)
	}
	if putBody.URL != "https://doc1.domain.com/" {
		t.Fatalf("url = %q, want https://doc1.domain.com/", putBody.URL)
	}

	conflictReq := signedAdminRequest(t, http.MethodPost, "/api/v1/documents/doc1", []byte("duplicate"), fixedTime, "ffff1111bbbb2222cccc3333dddd4444")
	conflictResp := httptest.NewRecorder()
	handler.ServeHTTP(conflictResp, conflictReq)
	if conflictResp.Code != http.StatusConflict || strings.TrimSpace(conflictResp.Body.String()) != `{"error":"already_exists"}` {
		t.Fatalf("create conflict response = (%d, %q)", conflictResp.Code, conflictResp.Body.String())
	}

	replaceReq := signedAdminRequest(t, http.MethodPut, "/api/v1/documents/doc1", []byte("hello world updated"), fixedTime, "abab1111bbbb2222cccc3333dddd4444")
	replaceResp := httptest.NewRecorder()
	handler.ServeHTTP(replaceResp, replaceReq)
	if replaceResp.Code != http.StatusOK {
		t.Fatalf("PUT status = %d, want %d, body=%s", replaceResp.Code, http.StatusOK, replaceResp.Body.String())
	}

	publicGetReq := httptest.NewRequest(http.MethodGet, "http://doc1.domain.com/", nil)
	publicGetReq.Host = "doc1.domain.com"
	publicGetResp := httptest.NewRecorder()
	handler.ServeHTTP(publicGetResp, publicGetReq)
	if publicGetResp.Code != http.StatusOK {
		t.Fatalf("public GET status = %d, want %d", publicGetResp.Code, http.StatusOK)
	}
	if body := publicGetResp.Body.String(); body != "hello world updated" {
		t.Fatalf("public GET body = %q, want %q", body, "hello world updated")
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
	if contentResp.Body.String() != "hello world updated" {
		t.Fatalf("content body = %q, want %q", contentResp.Body.String(), "hello world updated")
	}

	healthReq := httptest.NewRequest(http.MethodGet, "http://admin.domain.com/api/v1/health", nil)
	healthReq.Host = "admin.domain.com"
	healthResp := httptest.NewRecorder()
	handler.ServeHTTP(healthResp, healthReq)
	if healthResp.Code != http.StatusOK || strings.TrimSpace(healthResp.Body.String()) != `{"status":"ok"}` {
		t.Fatalf("health response = (%d, %q)", healthResp.Code, healthResp.Body.String())
	}

	authCheckReq := signedAdminRequest(t, http.MethodGet, "/api/v1/auth/check", nil, fixedTime, "adad1111bbbb2222cccc3333dddd4444")
	authCheckResp := httptest.NewRecorder()
	handler.ServeHTTP(authCheckResp, authCheckReq)
	if authCheckResp.Code != http.StatusOK || strings.TrimSpace(authCheckResp.Body.String()) != `{"status":"ok"}` {
		t.Fatalf("auth check response = (%d, %q)", authCheckResp.Code, authCheckResp.Body.String())
	}
	if cacheControl := authCheckResp.Header().Get("Cache-Control"); cacheControl != "no-store" {
		t.Fatalf("auth check Cache-Control = %q, want %q", cacheControl, "no-store")
	}

	missingReplaceReq := signedAdminRequest(t, http.MethodPut, "/api/v1/documents/missing", []byte("nope"), fixedTime, "acac1111bbbb2222cccc3333dddd4444")
	missingReplaceResp := httptest.NewRecorder()
	handler.ServeHTTP(missingReplaceResp, missingReplaceReq)
	if missingReplaceResp.Code != http.StatusNotFound || strings.TrimSpace(missingReplaceResp.Body.String()) != `{"error":"not_found"}` {
		t.Fatalf("replace missing response = (%d, %q)", missingReplaceResp.Code, missingReplaceResp.Body.String())
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

func TestAdminReplaySurvivesAppRestart(t *testing.T) {
	fixedTime := time.Date(2026, 4, 2, 12, 0, 0, 0, time.UTC)
	dbPath := filepath.Join(t.TempDir(), "data.db")
	handlerOne, cleanupOne := newTestHandlerWithDB(t, fixedTime, dbPath)

	reqOne := signedAdminRequest(t, http.MethodPost, "/api/v1/documents/doc1", []byte("hello world"), fixedTime, "ffff0000bbbb2222cccc3333dddd4444")
	respOne := httptest.NewRecorder()
	handlerOne.ServeHTTP(respOne, reqOne)
	if respOne.Code != http.StatusCreated {
		t.Fatalf("first POST status = %d, want %d, body=%s", respOne.Code, http.StatusCreated, respOne.Body.String())
	}
	cleanupOne()

	handlerTwo, cleanupTwo := newTestHandlerWithDB(t, fixedTime, dbPath)
	defer cleanupTwo()

	reqTwo := signedAdminRequest(t, http.MethodPost, "/api/v1/documents/doc2", []byte("hello again"), fixedTime, "ffff0000bbbb2222cccc3333dddd4444")
	respTwo := httptest.NewRecorder()
	handlerTwo.ServeHTTP(respTwo, reqTwo)
	if respTwo.Code != http.StatusUnauthorized || strings.TrimSpace(respTwo.Body.String()) != `{"error":"replayed_nonce"}` {
		t.Fatalf("replay after restart response = (%d, %q)", respTwo.Code, respTwo.Body.String())
	}
}

func newTestHandler(t *testing.T, fixedTime time.Time) (http.Handler, func()) {
	t.Helper()
	return newTestHandlerWithDB(t, fixedTime, filepath.Join(t.TempDir(), "data.db"))
}

func newTestHandlerWithDB(t *testing.T, fixedTime time.Time, dbPath string) (http.Handler, func()) {
	t.Helper()
	cfg := config.Config{
		ListenAddr:         "127.0.0.1:8080",
		BaseDomain:         "domain.com",
		AdminHost:          "admin.domain.com",
		PresharedKey:       "01234567890123456789012345678901",
		DBPath:             dbPath,
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
	req.Header.Set(authsig.HeaderTimestamp, strconv.FormatInt(timestamp.Unix(), 10))
	req.Header.Set(authsig.HeaderNonce, nonceValue)
	canonical := authsig.CanonicalRequest(req, authsig.SHA256Hex(body), req.Header.Get(authsig.HeaderTimestamp), nonceValue)
	req.Header.Set(authsig.HeaderSignature, hex.EncodeToString(authsig.Sign([]byte("01234567890123456789012345678901"), canonical)))
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	return req
}

func decodeJSON(t *testing.T, raw string, target any) {
	t.Helper()
	if err := json.Unmarshal([]byte(raw), target); err != nil {
		t.Fatalf("json.Unmarshal(%q) error = %v", raw, err)
	}
}
