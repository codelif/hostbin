package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/codelif/hostbin/internal/server/requestmeta"
)

const requestIDHeader = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := randomHex(12)
		meta := &requestmeta.Meta{RequestID: requestID}
		r = r.WithContext(requestmeta.WithContext(r.Context(), meta))

		w.Header().Set(requestIDHeader, requestID)
		next.ServeHTTP(w, r)
	})
}

func randomHex(size int) string {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "000000000000000000000000"
	}

	return hex.EncodeToString(b)
}
