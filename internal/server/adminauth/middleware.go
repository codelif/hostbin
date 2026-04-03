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

package adminauth

import (
	"bytes"
	"crypto/hmac"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/codelif/hostbin/internal/clock"
	"github.com/codelif/hostbin/internal/domain/hosts"
	"github.com/codelif/hostbin/internal/protocol/adminv1"
	"github.com/codelif/hostbin/internal/protocol/authsig"
	"github.com/codelif/hostbin/internal/server/nonce"
)

type Verifier struct {
	adminHost     string
	secret        []byte
	clock         clock.Clock
	timestampSkew time.Duration
	nonceStore    nonce.Store
}

func NewVerifier(adminHost string, secret []byte, clock clock.Clock, timestampSkew time.Duration, nonceStore nonce.Store) *Verifier {
	return &Verifier{
		adminHost:     adminHost,
		secret:        secret,
		clock:         clock,
		timestampSkew: timestampSkew,
		nonceStore:    nonceStore,
	}
}

func (v *Verifier) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		host, err := hosts.NormalizeHost(c.Request.Host)
		if err != nil || host != v.adminHost {
			abort(c, http.StatusUnauthorized, adminv1.ErrorUnauthorized)
			return
		}

		headers, errorCode := ParseHeaders(c.Request)
		if errorCode != "" {
			abort(c, statusForError(errorCode), errorCode)
			return
		}

		now := v.clock.Now().UTC()
		if !withinSkew(now, headers.Timestamp, v.timestampSkew) {
			abort(c, http.StatusUnauthorized, adminv1.ErrorInvalidTimestamp)
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			if isBodyTooLarge(err) {
				abort(c, http.StatusRequestEntityTooLarge, adminv1.ErrorDocumentTooLarge)
				return
			}

			abort(c, http.StatusBadRequest, adminv1.ErrorBadRequest)
			return
		}

		canonical := authsig.CanonicalRequest(c.Request, authsig.SHA256Hex(body), headers.TimestampRaw, headers.Nonce)
		expected := authsig.Sign(v.secret, canonical)
		if !hmac.Equal(expected, headers.Signature) {
			abort(c, http.StatusUnauthorized, adminv1.ErrorInvalidSignature)
			return
		}

		if err := v.nonceStore.UseOnce(headers.Nonce, now); err != nil {
			if errors.Is(err, nonce.ErrReplayed) {
				abort(c, http.StatusUnauthorized, adminv1.ErrorReplayedNonce)
				return
			}

			abort(c, http.StatusInternalServerError, adminv1.ErrorInternal)
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		c.Next()
	}
}

func withinSkew(now, timestamp time.Time, skew time.Duration) bool {
	delta := now.Sub(timestamp)
	if delta < 0 {
		delta = -delta
	}
	return delta <= skew
}

func statusForError(code string) int {
	switch code {
	case adminv1.ErrorInvalidTimestamp, adminv1.ErrorInvalidSignature, adminv1.ErrorReplayedNonce, adminv1.ErrorUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusBadRequest
	}
}

func abort(c *gin.Context, status int, code string) {
	c.AbortWithStatusJSON(status, adminv1.ErrorResponse{Error: code})
}

func isBodyTooLarge(err error) bool {
	var target *http.MaxBytesError
	return errors.As(err, &target)
}
