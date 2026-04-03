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

package hosts

import "testing"

func TestClassifyHost(t *testing.T) {
	reserved := map[string]struct{}{"admin": {}, "www": {}}

	tests := []struct {
		name     string
		host     string
		wantKind Kind
		wantSlug string
		wantErr  bool
	}{
		{name: "public host", host: "doc1.domain.com", wantKind: KindPublic, wantSlug: "doc1"},
		{name: "public host with port", host: "doc1.domain.com:443", wantKind: KindPublic, wantSlug: "doc1"},
		{name: "admin host", host: "admin.domain.com", wantKind: KindAdmin},
		{name: "base domain only", host: "domain.com", wantErr: true},
		{name: "multi level subdomain", host: "a.b.domain.com", wantErr: true},
		{name: "reserved public slug", host: "www.domain.com", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			info, err := ClassifyHost(tc.host, "domain.com", "admin.domain.com", reserved)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for host %q", tc.host)
				}
				return
			}

			if err != nil {
				t.Fatalf("ClassifyHost(%q) error = %v", tc.host, err)
			}
			if info.Kind != tc.wantKind {
				t.Fatalf("kind = %q, want %q", info.Kind, tc.wantKind)
			}
			if info.Slug != tc.wantSlug {
				t.Fatalf("slug = %q, want %q", info.Slug, tc.wantSlug)
			}
		})
	}
}
