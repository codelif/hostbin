package nonce

import (
	"errors"
	"time"
)

var ErrReplayed = errors.New("replayed nonce")

type Store interface {
	UseOnce(nonce string, now time.Time) error
	CleanupExpired(now time.Time)
}
