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
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/codelif/hostbin/internal/clock"
	"github.com/codelif/hostbin/internal/protocol/adminv1"
)

type rateLimitEntry struct {
	count   int
	resetAt time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	clock   clock.Clock
	limit   int
	window  time.Duration
	entries map[string]rateLimitEntry
}

func NewRateLimiter(appClock clock.Clock, limit int, window time.Duration) gin.HandlerFunc {
	limiter := &RateLimiter{
		clock:   appClock,
		limit:   limit,
		window:  window,
		entries: make(map[string]rateLimitEntry),
	}

	return limiter.Middleware()
}

func (l *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if l.allow(clientKey(c.Request.RemoteAddr)) {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusTooManyRequests, adminv1.ErrorResponse{Error: adminv1.ErrorRateLimited})
	}
}

func (l *RateLimiter) allow(key string) bool {
	now := l.clock.Now().UTC()

	l.mu.Lock()
	defer l.mu.Unlock()

	for entryKey, entry := range l.entries {
		if !entry.resetAt.After(now) {
			delete(l.entries, entryKey)
		}
	}

	entry, ok := l.entries[key]
	if !ok || !entry.resetAt.After(now) {
		l.entries[key] = rateLimitEntry{count: 1, resetAt: now.Add(l.window)}
		return true
	}

	if entry.count >= l.limit {
		return false
	}

	entry.count++
	l.entries[key] = entry
	return true
}

func clientKey(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil && host != "" {
		return host
	}

	if remoteAddr == "" {
		return "unknown"
	}

	return remoteAddr
}
