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
