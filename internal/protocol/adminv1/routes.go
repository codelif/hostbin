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

package adminv1

import "net/url"

const (
	BasePath               = "/api/v1"
	HealthPath             = BasePath + "/health"
	AuthCheckPath          = BasePath + "/auth/check"
	DocumentsCollection    = BasePath + "/documents"
	AuthCheckRelativePath  = "/auth/check"
	DocumentsRelativePath  = "/documents"
	DocumentPathPattern    = "/documents/:slug"
	DocumentContentPattern = "/documents/:slug/content"
)

func DocumentPath(slug string) string {
	return DocumentsCollection + "/" + url.PathEscape(slug)
}

func DocumentContentPath(slug string) string {
	return DocumentPath(slug) + "/content"
}
