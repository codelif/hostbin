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

package nonce

import (
	"sync"
	"time"
)

type MemoryStore struct {
	mu      sync.Mutex
	entries map[string]time.Time
	ttl     time.Duration
}

func NewMemoryStore(ttl time.Duration) *MemoryStore {
	return &MemoryStore{
		entries: make(map[string]time.Time),
		ttl:     ttl,
	}
}

func (s *MemoryStore) UseOnce(nonce string, now time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanupExpiredLocked(now)

	if expiresAt, ok := s.entries[nonce]; ok && expiresAt.After(now) {
		return ErrReplayed
	}

	s.entries[nonce] = now.UTC().Add(s.ttl)
	return nil
}

func (s *MemoryStore) CleanupExpired(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanupExpiredLocked(now)
}

func (s *MemoryStore) cleanupExpiredLocked(now time.Time) {
	for key, expiresAt := range s.entries {
		if !expiresAt.After(now) {
			delete(s.entries, key)
		}
	}
}
