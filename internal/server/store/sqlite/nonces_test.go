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

package sqlite

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/codelif/hostbin/internal/server/nonce"
)

func TestNonceStoreUseOnce(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "nonces.db"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer func() { _ = db.Close() }()

	store := NewNonceStore(db, 5*time.Minute)
	now := time.Date(2026, 4, 3, 12, 0, 0, 0, time.UTC)

	if err := store.UseOnce("nonce-1", now); err != nil {
		t.Fatalf("first UseOnce() error = %v", err)
	}
	if err := store.UseOnce("nonce-1", now.Add(time.Minute)); !errors.Is(err, nonce.ErrReplayed) {
		t.Fatalf("second UseOnce() error = %v, want ErrReplayed", err)
	}
	if err := store.UseOnce("nonce-1", now.Add(6*time.Minute)); err != nil {
		t.Fatalf("expired UseOnce() error = %v", err)
	}
}

func TestNonceStorePersistsAcrossInstances(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonces.db")
	now := time.Date(2026, 4, 3, 12, 0, 0, 0, time.UTC)

	db1, err := Open(path)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	store1 := NewNonceStore(db1, 5*time.Minute)
	if err := store1.UseOnce("nonce-2", now); err != nil {
		t.Fatalf("store1.UseOnce() error = %v", err)
	}
	if err := db1.Close(); err != nil {
		t.Fatalf("db1.Close() error = %v", err)
	}

	db2, err := Open(path)
	if err != nil {
		t.Fatalf("reopen Open() error = %v", err)
	}
	defer func() { _ = db2.Close() }()

	store2 := NewNonceStore(db2, 5*time.Minute)
	if err := store2.UseOnce("nonce-2", now.Add(time.Minute)); !errors.Is(err, nonce.ErrReplayed) {
		t.Fatalf("store2.UseOnce() error = %v, want ErrReplayed", err)
	}
}
