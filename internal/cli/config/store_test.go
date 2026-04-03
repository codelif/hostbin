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

package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestStoreLoadSave(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hbcli", "config.toml")
	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	if _, err := store.Load(); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Load() error = %v, want ErrNotFound", err)
	}

	config := DefaultFile()
	config.ServerURL = "https://admin.example.com"
	config.AuthKey = "01234567890123456789012345678901"

	if err := store.Save(config); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.ServerURL != config.ServerURL {
		t.Fatalf("ServerURL = %q, want %q", loaded.ServerURL, config.ServerURL)
	}
	if loaded.AuthKey != config.AuthKey {
		t.Fatalf("AuthKey = %q, want %q", loaded.AuthKey, config.AuthKey)
	}

	info, err := os.Stat(store.Path())
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("mode = %o, want 600", info.Mode().Perm())
	}
}

func TestStoreSaveAllowsPartialConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hbcli", "config.toml")
	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	partial := File{ServerURL: "https://admin.example.com"}
	if err := store.Save(partial); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.ServerURL != partial.ServerURL {
		t.Fatalf("ServerURL = %q, want %q", loaded.ServerURL, partial.ServerURL)
	}
	if loaded.AuthKey != "" {
		t.Fatalf("AuthKey = %q, want empty", loaded.AuthKey)
	}
}
