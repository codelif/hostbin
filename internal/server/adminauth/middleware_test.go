package adminauth

import (
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/codelif/hostbin/internal/clock"
	"github.com/codelif/hostbin/internal/protocol/authsig"
	"github.com/codelif/hostbin/internal/server/middleware"
	"github.com/codelif/hostbin/internal/server/nonce"
)

func TestVerifierValidSignature(t *testing.T) {
	engine, fixedTime := testEngine(t)
	req := signedRequest(t, http.MethodPut, "/api/v1/documents/doc1", []byte("hello world"), fixedTime, "abcd1234abcd1234abcd1234abcd1234")
	resp := httptest.NewRecorder()

	engine.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
	}
	if body := resp.Body.String(); body != "hello world" {
		t.Fatalf("body = %q, want %q", body, "hello world")
	}
}

func TestVerifierRejectsBodyMismatch(t *testing.T) {
	engine, fixedTime := testEngine(t)
	req := signedRequest(t, http.MethodPut, "/api/v1/documents/doc1", []byte("signed body"), fixedTime, "bbbb1234abcd1234abcd1234abcd1234")
	req.Body = io.NopCloser(strings.NewReader("different body"))
	req.ContentLength = int64(len("different body"))
	resp := httptest.NewRecorder()

	engine.ServeHTTP(resp, req)

	assertJSONError(t, resp, http.StatusUnauthorized, "invalid_signature")
}

func TestVerifierRejectsPathMismatch(t *testing.T) {
	engine, fixedTime := testEngine(t)
	req := signedRequest(t, http.MethodPut, "/api/v1/documents/other", []byte("hello"), fixedTime, "cccc1234abcd1234abcd1234abcd1234")
	req.URL.Path = "/api/v1/documents/doc1"
	req.URL.RawPath = ""
	resp := httptest.NewRecorder()

	engine.ServeHTTP(resp, req)

	assertJSONError(t, resp, http.StatusUnauthorized, "invalid_signature")
}

func TestVerifierRejectsOldTimestamp(t *testing.T) {
	engine, fixedTime := testEngine(t)
	req := signedRequest(t, http.MethodPut, "/api/v1/documents/doc1", []byte("hello"), fixedTime.Add(-2*time.Minute), "dddd1234abcd1234abcd1234abcd1234")
	resp := httptest.NewRecorder()

	engine.ServeHTTP(resp, req)

	assertJSONError(t, resp, http.StatusUnauthorized, "invalid_timestamp")
}

func TestVerifierRejectsReplayButAllowsValidRetryAfterInvalidSignature(t *testing.T) {
	engine, fixedTime := testEngine(t)
	nonceValue := "eeee1234abcd1234abcd1234abcd1234"

	invalidReq := signedRequest(t, http.MethodPut, "/api/v1/documents/doc1", []byte("hello"), fixedTime, nonceValue)
	invalidReq.Header.Set(authsig.HeaderSignature, strings.Repeat("0", 64))
	invalidResp := httptest.NewRecorder()
	engine.ServeHTTP(invalidResp, invalidReq)
	assertJSONError(t, invalidResp, http.StatusUnauthorized, "invalid_signature")

	validReq := signedRequest(t, http.MethodPut, "/api/v1/documents/doc1", []byte("hello"), fixedTime, nonceValue)
	validResp := httptest.NewRecorder()
	engine.ServeHTTP(validResp, validReq)
	if validResp.Code != http.StatusOK {
		t.Fatalf("status after invalid signature = %d, want %d", validResp.Code, http.StatusOK)
	}

	replayReq := signedRequest(t, http.MethodPut, "/api/v1/documents/doc1", []byte("hello"), fixedTime, nonceValue)
	replayResp := httptest.NewRecorder()
	engine.ServeHTTP(replayResp, replayReq)
	assertJSONError(t, replayResp, http.StatusUnauthorized, "replayed_nonce")
}

func testEngine(t *testing.T) (*gin.Engine, time.Time) {
	t.Helper()
	gin.SetMode(gin.ReleaseMode)
	fixedTime := time.Date(2026, 4, 2, 12, 0, 0, 0, time.UTC)
	verifier := NewVerifier(
		"admin.domain.com",
		[]byte("01234567890123456789012345678901"),
		clock.Fixed{Time: fixedTime},
		time.Minute,
		nonce.NewMemoryStore(5*time.Minute),
	)

	engine := gin.New()
	engine.UseRawPath = true
	engine.UnescapePathValues = false
	engine.Use(middleware.LimitBodyBytes(1024), verifier.Middleware())
	engine.PUT("/api/v1/documents/:slug", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		c.String(http.StatusOK, string(body))
	})

	return engine, fixedTime
}

func signedRequest(t *testing.T, method, path string, body []byte, timestamp time.Time, nonceValue string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(string(body)))
	req.Host = "admin.domain.com"
	req.URL.RawPath = path
	req.Header.Set(authsig.HeaderTimestamp, strconv.FormatInt(timestamp.Unix(), 10))
	req.Header.Set(authsig.HeaderNonce, nonceValue)
	canonical := authsig.CanonicalRequest(req, authsig.SHA256Hex(body), req.Header.Get(authsig.HeaderTimestamp), nonceValue)
	req.Header.Set(authsig.HeaderSignature, hex.EncodeToString(authsig.Sign([]byte("01234567890123456789012345678901"), canonical)))
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	req.ContentLength = int64(len(body))
	return req
}

func assertJSONError(t *testing.T, resp *httptest.ResponseRecorder, wantStatus int, wantError string) {
	t.Helper()
	if resp.Code != wantStatus {
		t.Fatalf("status = %d, want %d", resp.Code, wantStatus)
	}
	if body := resp.Body.String(); body != `{"error":"`+wantError+`"}` {
		t.Fatalf("body = %q, want %q", body, `{"error":"`+wantError+`"}`)
	}
}
