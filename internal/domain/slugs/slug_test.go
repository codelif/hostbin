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

package slugs

import "testing"

func TestValidate(t *testing.T) {
	reserved := map[string]struct{}{"admin": {}}

	tests := []struct {
		name    string
		value   string
		wantErr error
	}{
		{name: "simple", value: "doc1"},
		{name: "hyphen", value: "my-notes"},
		{name: "reserved", value: "admin", wantErr: ErrReserved},
		{name: "underscore", value: "bad_slug", wantErr: ErrInvalid},
		{name: "leading hyphen", value: "-bad", wantErr: ErrInvalid},
		{name: "trailing hyphen", value: "bad-", wantErr: ErrInvalid},
		{name: "uppercase", value: "Doc1", wantErr: ErrInvalid},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(tc.value, reserved)
			if err != tc.wantErr {
				t.Fatalf("Validate(%q) error = %v, want %v", tc.value, err, tc.wantErr)
			}
		})
	}
}
