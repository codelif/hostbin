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
	"encoding/hex"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/codelif/hostbin/internal/protocol/adminv1"
	"github.com/codelif/hostbin/internal/protocol/authsig"
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
