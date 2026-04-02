package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
)

func SHA256Hex(body []byte) string {
	hash := sha256.Sum256(body)
	return hex.EncodeToString(hash[:])
}

func CanonicalRequest(r *http.Request, bodyHash, timestamp, nonce string) string {
	parts := []string{
		strings.ToUpper(r.Method),
		rawPath(r),
		r.URL.RawQuery,
		bodyHash,
		timestamp,
		nonce,
	}

	return strings.Join(parts, "\n")
}

func Sign(secret []byte, canonical string) []byte {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(canonical))
	return mac.Sum(nil)
}

func rawPath(r *http.Request) string {
	if r.URL.RawPath != "" {
		return r.URL.RawPath
	}
	return r.URL.EscapedPath()
}
