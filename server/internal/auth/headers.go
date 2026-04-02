package auth

import (
	"encoding/hex"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	HeaderTimestamp = "X-Timestamp"
	HeaderNonce     = "X-Nonce"
	HeaderSignature = "X-Signature"
)

var (
	hexNoncePattern     = regexp.MustCompile(`^[0-9a-f]{32,128}$`)
	base64URLPattern    = regexp.MustCompile(`^[A-Za-z0-9_-]{22,128}$`)
	signatureHexPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

type Headers struct {
	TimestampRaw string
	Timestamp    time.Time
	Nonce        string
	Signature    []byte
}

func ParseHeaders(r *http.Request) (Headers, string) {
	timestampRaw := r.Header.Get(HeaderTimestamp)
	nonce := r.Header.Get(HeaderNonce)
	signatureRaw := r.Header.Get(HeaderSignature)

	if timestampRaw == "" || nonce == "" || signatureRaw == "" {
		return Headers{}, "unauthorized"
	}

	timestampUnix, err := strconv.ParseInt(timestampRaw, 10, 64)
	if err != nil {
		return Headers{}, "invalid_timestamp"
	}

	if !hexNoncePattern.MatchString(nonce) && !base64URLPattern.MatchString(nonce) {
		return Headers{}, "unauthorized"
	}

	if !signatureHexPattern.MatchString(signatureRaw) {
		return Headers{}, "invalid_signature"
	}

	signature, err := hex.DecodeString(signatureRaw)
	if err != nil {
		return Headers{}, "invalid_signature"
	}

	return Headers{
		TimestampRaw: timestampRaw,
		Timestamp:    time.Unix(timestampUnix, 0).UTC(),
		Nonce:        nonce,
		Signature:    signature,
	}, ""
}
