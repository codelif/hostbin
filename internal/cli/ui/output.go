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

package ui

import (
	"fmt"
	"io"

	cliconfig "github.com/codelif/hostbin/internal/cli/config"
	"github.com/codelif/hostbin/internal/cli/format"
	"github.com/codelif/hostbin/internal/protocol/adminv1"
)

func PrintConfigSummary(w io.Writer, path string, cfg cliconfig.File) error {
	normalized := cfg.Normalized()
	_, err := fmt.Fprintf(w,
		"hbcli configuration\n  path:       %s\n  server_url: %s\n  auth_key:   %s\n  timeout:    %s\n  editor:     %s\n  color:      %s\n",
		path,
		valueOrUnset(normalized.ServerURL),
		valueOrUnset(cliconfig.MaskSecret(normalized.AuthKey)),
		normalized.Timeout,
		valueOrUnset(normalized.Editor),
		normalized.Color,
	)
	return err
}

func valueOrUnset(value string) string {
	if value == "" {
		return "(unset)"
	}

	return value
}

func PrintDocumentSummary(w io.Writer, action string, doc adminv1.DocumentResponse) error {
	_, err := fmt.Fprintf(w,
		"%s %s\n  URL: %s\n  Size: %s\n  Updated: %s\n",
		action,
		doc.Slug,
		doc.URL,
		format.Bytes(doc.SizeBytes),
		format.Timestamp(doc.UpdatedAt),
	)
	return err
}

func PrintDocumentInfo(w io.Writer, doc adminv1.DocumentResponse) error {
	_, err := fmt.Fprintf(w,
		"%s\n  URL: %s\n  Size: %s\n  SHA256: %s\n  Created: %s\n  Updated: %s\n",
		doc.Slug,
		doc.URL,
		format.Bytes(doc.SizeBytes),
		doc.SHA256,
		doc.CreatedAt,
		doc.UpdatedAt,
	)
	return err
}

func PrintDeleteSummary(w io.Writer, slug string) error {
	_, err := fmt.Fprintf(w, "Deleted %s\n", slug)
	return err
}
