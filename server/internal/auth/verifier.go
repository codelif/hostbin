package auth

import (
	"bytes"
	"crypto/hmac"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"hostbin/internal/clock"
	"hostbin/internal/nonce"
	"hostbin/internal/router"
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
		host, err := router.NormalizeHost(c.Request.Host)
		if err != nil || host != v.adminHost {
			abort(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		headers, errorCode := ParseHeaders(c.Request)
		if errorCode != "" {
			abort(c, statusForError(errorCode), errorCode)
			return
		}

		now := v.clock.Now().UTC()
		if !withinSkew(now, headers.Timestamp, v.timestampSkew) {
			abort(c, http.StatusUnauthorized, "invalid_timestamp")
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			if isBodyTooLarge(err) {
				abort(c, http.StatusRequestEntityTooLarge, "document_too_large")
				return
			}

			abort(c, http.StatusBadRequest, "bad_request")
			return
		}

		canonical := CanonicalRequest(c.Request, SHA256Hex(body), headers.TimestampRaw, headers.Nonce)
		expected := Sign(v.secret, canonical)
		if !hmac.Equal(expected, headers.Signature) {
			abort(c, http.StatusUnauthorized, "invalid_signature")
			return
		}

		if err := v.nonceStore.UseOnce(headers.Nonce, now); err != nil {
			if errors.Is(err, nonce.ErrReplayed) {
				abort(c, http.StatusUnauthorized, "replayed_nonce")
				return
			}

			abort(c, http.StatusInternalServerError, "internal_error")
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
	case "invalid_timestamp", "invalid_signature", "replayed_nonce", "unauthorized":
		return http.StatusUnauthorized
	default:
		return http.StatusBadRequest
	}
}

func abort(c *gin.Context, status int, code string) {
	c.AbortWithStatusJSON(status, gin.H{"error": code})
}

func isBodyTooLarge(err error) bool {
	var target *http.MaxBytesError
	return errors.As(err, &target)
}
