package authsig

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

func SetSignedHeaders(r *http.Request, body []byte, secret []byte, now time.Time) error {
	nonce, err := NewNonce()
	if err != nil {
		return err
	}

	return SetSignedHeadersWithNonce(r, body, secret, now, nonce)
}

func SetSignedHeadersWithNonce(r *http.Request, body []byte, secret []byte, now time.Time, nonce string) error {
	if err := ValidateSharedSecret(string(secret)); err != nil {
		return err
	}

	timestamp := fmt.Sprintf("%d", now.UTC().Unix())
	canonical := CanonicalRequest(r, SHA256Hex(body), timestamp, nonce)
	signature := hex.EncodeToString(Sign(secret, canonical))

	r.Header.Set(HeaderTimestamp, timestamp)
	r.Header.Set(HeaderNonce, nonce)
	r.Header.Set(HeaderSignature, signature)

	return nil
}

func NewNonce() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}
