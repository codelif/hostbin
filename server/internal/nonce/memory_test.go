package nonce

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestMemoryStoreUseOnceIsAtomic(t *testing.T) {
	store := NewMemoryStore(5 * time.Minute)
	now := time.Date(2026, 4, 2, 12, 0, 0, 0, time.UTC)

	var successCount int32
	var replayCount int32

	start := make(chan struct{})
	var wg sync.WaitGroup
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			err := store.UseOnce("same-nonce", now)
			switch err {
			case nil:
				atomic.AddInt32(&successCount, 1)
			case ErrReplayed:
				atomic.AddInt32(&replayCount, 1)
			default:
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}

	close(start)
	wg.Wait()

	if successCount != 1 {
		t.Fatalf("successCount = %d, want 1", successCount)
	}
	if replayCount != 15 {
		t.Fatalf("replayCount = %d, want 15", replayCount)
	}
}
