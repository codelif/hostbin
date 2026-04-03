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
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/codelif/hostbin/internal/clock"
)

func TestRateLimiterBlocksRepeatedRequests(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(NewRateLimiter(clock.Fixed{Time: time.Date(2026, 4, 3, 12, 0, 0, 0, time.UTC)}, 2, time.Minute))
	engine.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "198.51.100.10:1234"
		resp := httptest.NewRecorder()
		engine.ServeHTTP(resp, req)
		if resp.Code != http.StatusOK {
			t.Fatalf("request %d status = %d, want %d", i+1, resp.Code, http.StatusOK)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "198.51.100.10:1234"
	resp := httptest.NewRecorder()
	engine.ServeHTTP(resp, req)
	if resp.Code != http.StatusTooManyRequests {
		t.Fatalf("blocked status = %d, want %d", resp.Code, http.StatusTooManyRequests)
	}
	if body := resp.Body.String(); body != `{"error":"rate_limited"}` {
		t.Fatalf("blocked body = %q, want %q", body, `{"error":"rate_limited"}`)
	}
}

func TestRateLimiterSeparatesClients(t *testing.T) {
	limiter := &RateLimiter{
		clock:   clock.Fixed{Time: time.Date(2026, 4, 3, 12, 0, 0, 0, time.UTC)},
		limit:   1,
		window:  time.Minute,
		entries: make(map[string]rateLimitEntry),
	}

	if !limiter.allow("198.51.100.10") {
		t.Fatal("first client should be allowed")
	}
	if !limiter.allow("198.51.100.11") {
		t.Fatal("second client should be allowed")
	}
	if limiter.allow("198.51.100.10") {
		t.Fatal("repeat request should be blocked")
	}
}
