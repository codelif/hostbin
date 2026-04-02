package adminauth

import (
	"encoding/hex"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"hostbin/internal/protocol/adminv1"
	"hostbin/internal/protocol/authsig"
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
	timestampRaw := r.Header.Get(authsig.HeaderTimestamp)
	nonce := r.Header.Get(authsig.HeaderNonce)
	signatureRaw := r.Header.Get(authsig.HeaderSignature)

	if timestampRaw == "" || nonce == "" || signatureRaw == "" {
		return Headers{}, adminv1.ErrorUnauthorized
	}

	timestampUnix, err := strconv.ParseInt(timestampRaw, 10, 64)
	if err != nil {
		return Headers{}, adminv1.ErrorInvalidTimestamp
	}

	if !hexNoncePattern.MatchString(nonce) && !base64URLPattern.MatchString(nonce) {
		return Headers{}, adminv1.ErrorUnauthorized
	}

	if !signatureHexPattern.MatchString(signatureRaw) {
		return Headers{}, adminv1.ErrorInvalidSignature
	}

	signature, err := hex.DecodeString(signatureRaw)
	if err != nil {
		return Headers{}, adminv1.ErrorInvalidSignature
	}

	return Headers{
		TimestampRaw: timestampRaw,
		Timestamp:    time.Unix(timestampUnix, 0).UTC(),
		Nonce:        nonce,
		Signature:    signature,
	}, ""
}
